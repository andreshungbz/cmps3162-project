package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andreshungbz/cmps3162-project/internal/validator"
)

var (
	ErrDuplicatePassport = errors.New("duplicate passport")
)

// Guest maps the guest entity (subtype of the person entity).
type Guest struct {
	// person attributes
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	// guest attributes
	ID             int64  `json:"-"`
	PassportNumber string `json:"passport_number"`
	ContactEmail   string `json:"contact_email"`
	ContactPhone   string `json:"contact_phone"`
}

// ValidateContactEmail validates an email against the regular expression.
func ValidateContactEmail(v *validator.Validator, email string) {
	v.Check(email != "", "contact_email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "contact_email", "must be a valid contact email address")
}

// ValidateGuest performs validation checks for a guest record.
func ValidateGuest(v *validator.Validator, guest *Guest) {
	// person attributes
	v.Check(guest.Name != "", "name", "must be provided")
	v.Check(guest.Gender != "", "gender", "must be provided")
	v.Check(guest.Street != "", "street", "must be provided")
	v.Check(guest.City != "", "city", "must be provided")
	v.Check(guest.Country != "", "country", "must be provided")

	// guest attributes
	v.Check(guest.PassportNumber != "", "passport_number", "must be provided")
	ValidateContactEmail(v, guest.ContactEmail)
	v.Check(guest.ContactPhone != "", "contact_phone", "must be provided")
}

// GuestModel holds the database handler.
type GuestModel struct {
	DB *sql.DB
}

// Insert creates a guest record (person + guest tables).
func (g GuestModel) Insert(guest *Guest) error {
	query := `SELECT * FROM fn_create_guest($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		// person attributes
		guest.Name,
		guest.Gender,
		guest.Street,
		guest.City,
		guest.Country,
		// guest attributes
		guest.PassportNumber,
		guest.ContactEmail,
		guest.ContactPhone,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := g.DB.QueryRowContext(ctx, query, args...).Scan(&guest.ID, &guest.CreatedAt, &guest.ModifiedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "guest_passport_number_key" (23505)`:
			return ErrDuplicatePassport
		default:
			return err
		}
	}

	return nil
}

// GetByPassport retrieves a single guest record by passport_number.
func (g GuestModel) GetByPassport(passport string) (*Guest, error) {
	query := `SELECT * FROM fn_get_guest_by_passport($1)`
	var guest Guest

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := g.DB.QueryRowContext(ctx, query, passport).Scan(
		// person attributes
		&guest.Name,
		&guest.Gender,
		&guest.Street,
		&guest.City,
		&guest.Country,
		&guest.CreatedAt,
		&guest.ModifiedAt,
		// guest attributes
		&guest.ID,
		&guest.PassportNumber,
		&guest.ContactEmail,
		&guest.ContactPhone,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &guest, nil
}

// GetAll retrieves multiple guest records (filterable).
func (g GuestModel) GetAll(name string, country string, filters Filters) ([]*Guest, Metadata, error) {
	// PostgreSQL full-text search notes
	// - to_tsvector in simple configuration breaks string to lower case lexemes.
	// - plainto_tsquery in simple configuration normalizes the query term.
	//		e.g. "John Smith" -> 'john' + 'smith'
	// - @@ is the matches operator.
	query := fmt.Sprintf(`
		SELECT *
		FROM fn_get_guests($1,$2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`,
		filters.sortColumn(),
		filters.sortDirection(),
	)

	// limit is used for page_size, and offset is used for page
	args := []any{name, country, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// retrieves rows from the database
	rows, err := g.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	// construct the array of guests
	totalRecords := 0
	guests := []*Guest{}
	for rows.Next() {
		var guest Guest

		err := rows.Scan(
			&totalRecords,
			// person attributes
			&guest.Name,
			&guest.Gender,
			&guest.Street,
			&guest.City,
			&guest.Country,
			&guest.CreatedAt,
			&guest.ModifiedAt,
			// guest attributes
			&guest.ID,
			&guest.PassportNumber,
			&guest.ContactEmail,
			&guest.ContactPhone,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		guests = append(guests, &guest)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// construct Metadata object
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return guests, metadata, nil
}

// Update modifies a guest record by passport_number (person + guest).
func (g GuestModel) Update(guest *Guest) error {
	query := `SELECT fn_update_guest($1, $2, $3, $4, $5, $6, $7, $8)`

	args := []any{
		// person attributes
		guest.Name,
		guest.Gender,
		guest.Street,
		guest.City,
		guest.Country,
		// guest attributes
		guest.PassportNumber,
		guest.ContactEmail,
		guest.ContactPhone,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "guest_passport_number_key" (23505)`:
			return ErrDuplicatePassport
		default:
			return err
		}
	}

	// No rows being affected means the guest's passport is not in the database

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Delete removes a guest record by passport_number (cascades to reservation and registration).
func (g GuestModel) Delete(passport string) error {
	query := `SELECT fn_delete_guest($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, passport)
	if err != nil {
		return err
	}

	// No rows being affected means the guest's passport is not in the database

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
