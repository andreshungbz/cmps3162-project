package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/cmps3162-project/internal/data"
	"github.com/andreshungbz/cmps3162-project/internal/validator"
)

// createEmployeeHandler calls Employee.Insert.
// Writes JSON of the created employee record.
func (app *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		// person attributes
		Name    string `json:"name"`
		Gender  string `json:"gender"`
		Street  string `json:"street"`
		City    string `json:"city"`
		Country string `json:"country"`
		// employee attributes
		HotelID    int64   `json:"hotel_id"`
		Department string  `json:"department"`
		ManagerID  *int64  `json:"manager_id,omitempty"`
		Salary     float64 `json:"salary"`
		SSN        string  `json:"ssn"`
		WorkEmail  string  `json:"work_email"`
		WorkPhone  string  `json:"work_phone"`
		Password   string  `json:"password"`
		// role-specific attributes
		Role       string  `json:"role"`
		HotelOwner *bool   `json:"hotel_owner,omitempty"`
		Shift      *string `json:"shift,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	employee := &data.Employee{
		// person attributes
		Name:    input.Name,
		Gender:  input.Gender,
		Street:  input.Street,
		City:    input.City,
		Country: input.Country,
		// employee attributes
		HotelID:    input.HotelID,
		Department: input.Department,
		ManagerID:  input.ManagerID,
		Salary:     input.Salary,
		SSN:        input.SSN,
		WorkEmail:  input.WorkEmail,
		WorkPhone:  input.WorkPhone,
		Employed:   true,
		Activated:  false,
		// role-specific attributes
		Role:       input.Role,
		HotelOwner: input.HotelOwner,
		Shift:      input.Shift,
	}

	// generate password hash
	err = employee.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateEmployee(v, employee); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Employee.Insert(employee)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSSN):
			v.AddError("ssn", "an employee with this social security number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("work_email", "an employee with this work email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicatePhone):
			v.AddError("work_phone", "an employee with this work phone number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	// err = app.models.Permissions.AddForEmployee(employee.ID, "healthcheck:read")
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	err = app.writeJSON(w, http.StatusCreated, envelope{"employee": employee}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showEmployeeHandler calls Employee.GetByEmail.
// Writes JSON of the retrieved employee record.
func (app *application) showEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	email := app.readStringParam("email", r)
	if email == "" {
		app.notFoundResponse(w, r)
		return
	}

	employee, err := app.models.Employee.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employee": employee}, nil)
}

// listEmployeesHandler calls Employee.GetAll.
// Writes JSON of the list of filtered employee records and a metadata object.
func (app *application) listEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
		Role string
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Role = app.readURLString(qs, "role", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readURLString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "name", "department",
		"-id", "-name", "-department",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	employees, metadata, err := app.models.Employee.GetAll(input.Role, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employees": employees, "metadata": metadata}, nil)
}

// updateEmployeeHandler calls Employee.Update.
// Writes JSON of the updated employee record.
func (app *application) updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	email := app.readStringParam("email", r)
	if email == "" {
		app.notFoundResponse(w, r)
		return
	}

	employee, err := app.models.Employee.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		// person attributes
		Name    *string `json:"name"`
		Gender  *string `json:"gender"`
		Street  *string `json:"street"`
		City    *string `json:"city"`
		Country *string `json:"country"`
		// employee attributes
		HotelID    *int64   `json:"hotel_id,omitempty"`
		Department *string  `json:"department"`
		ManagerID  *int64   `json:"manager_id,omitempty"`
		Salary     *float64 `json:"salary"`
		SSN        *string  `json:"ssn"`
		WorkEmail  *string  `json:"work_email"`
		WorkPhone  *string  `json:"work_phone"`
		Password   *string  `json:"password"`
		Employed   *bool    `json:"employed"`
		Activated  *bool    `json:"activated"`
		// role-specific attributes
		Role       *string `json:"role"`
		HotelOwner *bool   `json:"hotel_owner,omitempty"`
		Shift      *string `json:"shift,omitempty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// person attributes
	if input.Name != nil {
		employee.Name = *input.Name
	}
	if input.Gender != nil {
		employee.Gender = *input.Gender
	}
	if input.Street != nil {
		employee.Street = *input.Street
	}
	if input.City != nil {
		employee.City = *input.City
	}
	if input.Country != nil {
		employee.Country = *input.Country
	}

	// employee attributes
	if input.HotelID != nil {
		employee.HotelID = *input.HotelID
	}
	if input.Department != nil {
		employee.Department = *input.Department
	}
	if input.ManagerID != nil {
		employee.ManagerID = input.ManagerID
	}
	if input.Salary != nil {
		employee.Salary = *input.Salary
	}
	if input.SSN != nil {
		employee.SSN = *input.SSN
	}
	if input.WorkEmail != nil {
		employee.WorkEmail = *input.WorkEmail
	}
	if input.WorkPhone != nil {
		employee.WorkPhone = *input.WorkPhone
	}

	// generate password hash
	if input.Password != nil {
		err = employee.Password.Set(*input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if input.Employed != nil {
		employee.Employed = *input.Employed
	}
	if input.Activated != nil {
		employee.Activated = *input.Activated
	}

	// role-specific attributes
	if input.Role != nil {
		employee.Role = *input.Role
	}
	if input.HotelOwner != nil {
		employee.HotelOwner = input.HotelOwner
	}
	if input.Shift != nil {
		employee.Shift = input.Shift
	}

	v := validator.New()
	if data.ValidateEmployee(v, employee); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Employee.Update(employee)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSSN):
			v.AddError("ssn", "an employee with this social security number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("work_email", "an employee with this work email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicatePhone):
			v.AddError("work_phone", "an employee with this work phone number already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"employee": employee}, nil)
}

// deleteEmployeeHandler calls Employee.Delete.
// Writes JSON of a successful deletion message.
func (app *application) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readInt64Param("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Employee.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "employee successfully deleted"}, nil)
}

// activateEmployeeHandler calls the necessary models to activate an employee.
// Writes JSON of the updated employee record.
func (app *application) activateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate token
	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// get employee by token
	employee, err := app.models.Employee.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	// set activated to true
	employee.Activated = true

	// update employee record
	err = app.models.Employee.Update(employee)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	// delete all activation tokens the for the employee after successful activation
	err = app.models.Tokens.DeleteAllForPerson(data.ScopeActivation, employee.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"employee": employee}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
