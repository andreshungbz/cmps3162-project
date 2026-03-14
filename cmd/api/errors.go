package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// logError writes server-side error messages. It records the error
// and HTTP request method and URI.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

// errorResponse writes error messages to the client in JSON.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse sends a generic error message in JSON indicating
// a problem with the server.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFound Response sends a 404 HTTP status code.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends a 405 HTTP status code.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse sends a 400 HTTP status code, typically for invalid JSON.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse sends a 422 HTTP status code, typically for
// validly structured JSON whose field values are not valid.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// editConflictResponse sends a 409 HTTP status code for concurrency conflicts
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "Unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// rateLimitExceededResponse sends a 429 HTTP status code for too many requests.
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retry time.Duration) {
	message := "Rate limit exceeded"

	// set the Retry-After header in seconds format, rounding up
	w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(retry.Seconds()))))

	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// invalidCredentialsResponse sends a 401 HTTP status code for unauthorized.
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// invalidAuthenticationTokenResponse sends a 401 HTTP status code for unauthorized.
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Www-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// authenticationRequiredResponse sends a 401 HTTP status code for unauthorized.
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Www-Authenticate", "Bearer")

	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// inactiveAccountResponse sends a 403 HTTP status code for forbidden.
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your employee account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// notPermittedResponse sends a 403 HTTP status code for forbidden.
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your employee account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
