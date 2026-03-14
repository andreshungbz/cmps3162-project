package main

import (
	"compress/gzip"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

// recoverPanic ensures that in the case of a panic, a Connection header of
// 'close' is sent to the client.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requestLogger logs the HTTP request's method and URL path.
func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("Request received", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// rateLimit uses a client's IP address to limit their rate.
func (app *application) rateLimit(next http.Handler) http.Handler {
	if !app.config.limiter.enabled {
		return next
	}

	// client holds a rate limiter and the last seen time.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// goroutine for removing stale entries from the list of clients
	go func() {
		time.Sleep(time.Minute)

		mu.Lock() // ensure no concurrency conflicts
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get client ip address
		ip := realip.FromRequest(r)

		mu.Lock()

		// add client if not in the list already
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
		}

		// update client's last seen time
		clients[ip].lastSeen = time.Now()

		// if client exceeds the rate limit, send an error response
		if !clients[ip].limiter.Allow() {
			// claim the next future token, determine its delay, then cancel it
			res := clients[ip].limiter.Reserve()
			delay := res.Delay()
			res.Cancel()

			mu.Unlock()
			app.rateLimitExceededResponse(w, r, delay)
			return
		}

		mu.Unlock()

		// if client is not rate-limited, call next handler
		next.ServeHTTP(w, r)
	})
}

// enableCORS configures browser CORS by reflecting a request's origin if they are in
// the list of trusted origins configured on server start. It also handles CORS
// preflight requests appropriately.
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// indicator for caches that these responses may vary
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		// retrieve the Origin header of the request
		origin := r.Header.Get("Origin")

		if origin != "" {
			// loop through every configured trusted origin
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] { // on match
					// set that origin for Access-Control-Allow-Origin,
					// allowing cross-origin requests
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// handle CORS preflight requests by checking OPTIONS and Access-Control-Request-Method
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// set the non-CORS-safe HTTP methods
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// write 200 instead of 204 No Content for browser compatibility
						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// metricsResponseWriter is a light wrapper around http.ResponseWriter that records
// status code and whether the header was written or not.
type metricsResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

// newMetricsResponseWriter is a factory function that sets the default status code to 200.
func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

// Header passes the http.ResponseWriter Header method.
func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

// WriteHeader passes the http.ResponseWriter WriteHeader method, writing
// the status code and indicating the header was written.
func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

// Write passes the http.ResponseWriter Write method, setting the header
// written to true.
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

// Unwrap returns the existing http.ResponseWriter.
func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}

// metrics records statistics about requests and responses handled by the server.
func (app *application) metrics(next http.Handler) http.Handler {
	// NOTE: the total responses will always be 1 less than the total requests because
	// during the first request, the metrics are recorded before the first response is
	// sent.

	var (
		// REQ: total requests
		totalRequestsReceived = expvar.NewInt("total_requests_received")

		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")

		// REQ: total requests per route
		totalRequestsByRoute = expvar.NewMap("total_requests_by_route")

		// REQ: for error counts, consider any response with status code >= 400
		totalResponsesErrorCount = expvar.NewInt("total_responses_error_count")

		// REQ: average latency can be calculated by dividing the total processing time
		// by the total requests received
		averageLatency = expvar.NewFloat("average_latency_μs")
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing

		start := time.Now()
		mw := newMetricsResponseWriter(w)

		totalRequestsReceived.Add(1)
		totalRequestsByRoute.Add(r.URL.Path, 1)

		// Processing
		next.ServeHTTP(mw, r)

		// Post-processing

		totalResponsesSent.Add(1)
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)
		if mw.statusCode >= 400 {
			totalResponsesErrorCount.Add(1)
		}

		// record processing time
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)

		// update average latency
		averageLatency.Set(float64(totalProcessingTimeMicroseconds.Value()) / float64(totalRequestsReceived.Value()))
	})
}

// gzipResponseWriter is a light wrapper around http.ResponseWriter that
// compresses responses written in the gzip format.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write uses a *gzip.Writer instead of the http.ResponseWriter.
func (gzw gzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.writer.Write(b)
}

// gzip compress responses if the client accepts the gzip encoding in its HTTP request.
func (app *application) gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check that the client set the Accept-Encoding header to "gzip"
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// set Content-Encoding and account for caching
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		// create a new *gzip.Writer and gzipResponseWriter
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{
			ResponseWriter: w,
			writer:         gz,
		}

		// use the gzipResponseWriter in the next handler
		next.ServeHTTP(gzw, r)
	})
}
