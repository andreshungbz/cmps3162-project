package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andreshungbz/cmps3162-project/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateSSN   = errors.New("duplicate ssn")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrDuplicatePhone = errors.New("duplicate phone")
)

// Employee maps the employee entity (subtype of person).
type Employee struct {
	// person attributes
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	// employee attributes
	ID         int64    `json:"id"`
	HotelID    int64    `json:"hotel_id"`
	Department string   `json:"department"`
	ManagerID  *int64   `json:"manager_id,omitempty"`
	Salary     float64  `json:"salary"`
	SSN        string   `json:"ssn"`
	WorkEmail  string   `json:"work_email"`
	WorkPhone  string   `json:"work_phone"`
	Password   password `json:"-"`
	Employed   bool     `json:"employed"`
	Activated  bool     `json:"activated"`
	// role attributes
	Role       string  `json:"role"`
	HotelOwner *bool   `json:"hotel_owner,omitempty"`
	Shift      *string `json:"shift,omitempty"`
}

// password holds the plaintext and hash of a password.
type password struct {
	plaintext *string
	hash      []byte
}

// Set sets the calculated hash of the plaintext password, which is also set.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Matches compares the hash of a passed in password against the hash in the struct.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// ValidateWorkEmail validates an email against the regular expression.
func ValidateWorkEmail(v *validator.Validator, email string) {
	v.Check(email != "", "work_email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "work_email", "must be a valid work email address")
}

// ValidatePasswordPlaintext validates password constraints.
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// ValidateEmployee performs validation checks for an employee record.
func ValidateEmployee(v *validator.Validator, e *Employee) {
	// person attributes
	v.Check(e.Name != "", "name", "must be provided")
	v.Check(e.Gender != "", "gender", "must be provided")
	v.Check(e.Street != "", "street", "must be provided")
	v.Check(e.City != "", "city", "must be provided")
	v.Check(e.Country != "", "country", "must be provided")

	// employee attributes
	v.Check(e.HotelID > 0, "hotel_id", "must be provided")
	v.Check(e.Department != "", "department", "must be provided")
	v.Check(e.Salary >= 0, "salary", "must be provided")
	v.Check(e.SSN != "", "ssn", "must be provided")
	ValidateWorkEmail(v, e.WorkEmail)
	v.Check(e.WorkPhone != "", "work_phone", "must be provided")

	if e.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *e.Password.plaintext)
	}
	if e.Password.hash == nil {
		panic("missing password hash for user")
	}

	// role-specific attributes
	v.Check(e.Role != "", "role", "must be provided")

	switch e.Role {
	case "operations_manager":
		v.Check(e.HotelOwner != nil, "hotel_owner", "must be provided for operations_manager")
	case "front_desk", "housekeeper":
		v.Check(e.Shift != nil && *e.Shift != "", "shift", "must be provided for front_desk or housekeeper")
	}
}

// EmployeeModel holds the database handler.
type EmployeeModel struct {
	DB *sql.DB
}

// Insert creates a new employee (person + employee + role-specific).
func (m EmployeeModel) Insert(e *Employee) error {
	query := `SELECT * FROM fn_create_employee(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
	)`

	args := []any{
		// person attributes
		e.Name,
		e.Gender,
		e.Street,
		e.City,
		e.Country,
		// employee attributes
		e.HotelID,
		e.Department,
		e.ManagerID,
		e.Salary,
		e.SSN,
		e.WorkEmail,
		e.WorkPhone,
		e.Password.hash,
		e.Employed,
		e.Activated,
		// role-specific attributes
		e.Role,
		e.HotelOwner,
		e.Shift,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&e.ID, &e.CreatedAt, &e.ModifiedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_ssn_key" (23505)`:
			return ErrDuplicateSSN
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_work_email_key" (23505)`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_work_phone_key" (23505)`:
			return ErrDuplicatePhone
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a single employee record by work email.
func (m EmployeeModel) GetByEmail(email string) (*Employee, error) {
	query := `SELECT * FROM fn_get_employee_by_email($1)`
	var employee Employee

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		// person attributes
		&employee.Name,
		&employee.Gender,
		&employee.Street,
		&employee.City,
		&employee.Country,
		&employee.CreatedAt,
		&employee.ModifiedAt,
		// employee attributes
		&employee.ID,
		&employee.HotelID,
		&employee.Department,
		&employee.ManagerID,
		&employee.Salary,
		&employee.SSN,
		&employee.WorkEmail,
		&employee.WorkPhone,
		&employee.Password.hash,
		&employee.Employed,
		&employee.Activated,
		// role-specific attributes
		&employee.Role,
		&employee.HotelOwner,
		&employee.Shift,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
		case strings.Contains(err.Error(), "[employee-not-found]"):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &employee, nil
}

// GetAll retrieves multiple employee records with optional role filter and pagination.
func (m EmployeeModel) GetAll(role string, filters Filters) ([]*Employee, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM fn_get_employees($1)
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`,
		filters.sortColumn(),
		filters.sortDirection(),
	)

	args := []any{role, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	employees := []*Employee{}
	for rows.Next() {
		var employee Employee
		err := rows.Scan(
			&totalRecords,
			// person attributes
			&employee.Name,
			&employee.Gender,
			&employee.Street,
			&employee.City,
			&employee.Country,
			&employee.CreatedAt,
			&employee.ModifiedAt,
			// employee attributes
			&employee.ID,
			&employee.HotelID,
			&employee.Department,
			&employee.ManagerID,
			&employee.Salary,
			&employee.SSN,
			&employee.WorkEmail,
			&employee.WorkPhone,
			&employee.Password.hash,
			&employee.Employed,
			&employee.Activated,
			// role-specific attributes
			&employee.Role,
			&employee.HotelOwner,
			&employee.Shift,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		employees = append(employees, &employee)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return employees, metadata, nil
}

// Update modifies an existing employee record (person + employee + role-specific).
func (m EmployeeModel) Update(employee *Employee) error {
	query := `
	SELECT fn_update_employee(
    	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
	)`

	args := []any{
		// person attributes
		employee.Name,
		employee.Gender,
		employee.Street,
		employee.City,
		employee.Country,
		// employee attributes
		employee.ID,
		employee.HotelID,
		employee.Department,
		employee.ManagerID,
		employee.Salary,
		employee.SSN,
		employee.WorkEmail,
		employee.WorkPhone,
		employee.Password.hash,
		employee.Employed,
		employee.Activated,
		// role-specific attributes
		employee.Role,
		employee.HotelOwner,
		employee.Shift,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_ssn_key" (23505)`:
			return ErrDuplicateSSN
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_work_email_key" (23505)`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "employee_work_phone_key" (23505)`:
			return ErrDuplicatePhone
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes an employee record by id.
func (m EmployeeModel) Delete(id int64) error {
	query := `
		DELETE FROM person
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
