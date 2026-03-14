package main

import (
	"context"
	"net/http"

	"github.com/andreshungbz/cmps3162-project/internal/data"
)

// setup a custom type for employee to implement context in http.Request
type contextKey string

const userContextKey = contextKey("employee")

// contextSetEmployee returns a copy of the http.Request with employee set as context.
func (app *application) contextSetEmployee(r *http.Request, employee *data.Employee) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, employee)
	return r.WithContext(ctx)
}

// contextGetEmployee retrieves an Employee struct from the http.Request context.
func (app *application) contextGetEmployee(r *http.Request) *data.Employee {
	employee, ok := r.Context().Value(userContextKey).(*data.Employee)
	if !ok {
		panic("missing employee value in request context")
	}

	return employee
}
