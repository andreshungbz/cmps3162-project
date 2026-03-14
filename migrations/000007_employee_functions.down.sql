-- migrations/employee_functions.down.sql
-- Drops the employee functions.

DROP FUNCTION IF EXISTS fn_delete_employee(BIGINT);
DROP FUNCTION IF EXISTS fn_update_employee(
    TEXT, TEXT, TEXT, TEXT, TEXT, BIGINT, BIGINT, TEXT, BIGINT, NUMERIC, CITEXT, TEXT, TEXT, BOOLEAN, TEXT
);
DROP FUNCTION IF EXISTS fn_get_employee(BIGINT);
DROP FUNCTION IF EXISTS fn_create_employee(
    TEXT, TEXT, TEXT, TEXT, TEXT, BIGINT, TEXT, BIGINT, NUMERIC, TEXT, CITEXT, TEXT, BYTEA, TEXT, BOOLEAN, TEXT
);
