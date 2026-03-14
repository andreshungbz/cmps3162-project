package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/andreshungbz/cmps3162-project/internal/data"
	"github.com/andreshungbz/cmps3162-project/internal/validator"
)

// createActivationTokenHandler calls Tokens.New.
// Writes JSON of the created activation token.
func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkEmail string `json:"work_email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateWorkEmail(v, input.WorkEmail); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	employee, err := app.models.Employee.GetByEmail(input.WorkEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("work_email", "no matching work email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	// already activated
	if employee.Activated {
		v.AddError("work_email", "user has already been activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// create activation token
	token, err := app.models.Tokens.New(employee.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// NOTE: would send an email at this point but we will return it in the response

	err = app.writeJSON(w, http.StatusAccepted, envelope{"activation_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createAuthenticationTokenHandler exchanges an employee's email and password for an
// authentication token.
// Writes JSON of the authentication token.
func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkEmail string `json:"work_email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate email and password
	v := validator.New()
	data.ValidateWorkEmail(v, input.WorkEmail)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// retrieve employee by their email
	employee, err := app.models.Employee.GetByEmail(input.WorkEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// check if provided password matches
	match, err := employee.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// create authentication token
	token, err := app.models.Tokens.New(employee.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
