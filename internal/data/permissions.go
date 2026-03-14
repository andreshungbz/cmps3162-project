package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Permissions holds permission codes.
type Permissions []string

// Include checks whether a code exists in a Permissions object.
func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

// PermissionModel holds the database handler.
type PermissionModel struct {
	DB *sql.DB
}

// GetAllForEmployee returns all permission codes for an employee by ID.
func (m PermissionModel) GetAllForEmployee(employeeID int64) (Permissions, error) {
	query := `
		SELECT permission.code
		FROM permission
		INNER JOIN employee_permission ON employee_permission.permission_id = permission.id
		INNER JOIN employee ON employee_permission.employee_id = employee.id
		WHERE employee.id = $1`
	var permissions Permissions

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AddForEmployee adds permission codes to an employee.
func (m PermissionModel) AddForEmployee(employeeID int64, codes ...string) error {
	query := `
		INSERT INTO employee_permission
		SELECT $1, permission.id FROM permission WHERE permission.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, employeeID, pq.Array(codes))
	return err
}
