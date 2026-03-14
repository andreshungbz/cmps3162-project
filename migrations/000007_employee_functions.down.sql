-- migrations/employee_functions.down.sql
-- Drops the employee functions.

DROP FUNCTION IF EXISTS fn_update_employee(
    TEXT, TEXT, TEXT, TEXT, TEXT,
    BIGINT, BIGINT, TEXT, BIGINT, NUMERIC, TEXT, CITEXT, TEXT, BYTEA, BOOLEAN, BOOLEAN,
    TEXT, BOOLEAN, TEXT
);

DROP FUNCTION IF EXISTS fn_get_employees(TEXT);

DROP FUNCTION IF EXISTS fn_get_employee_for_token(BYTEA, TEXT);

DROP FUNCTION IF EXISTS fn_get_employee_by_email(CITEXT);

DROP FUNCTION IF EXISTS fn_create_employee(
    TEXT, TEXT, TEXT, TEXT, TEXT,
    BIGINT, TEXT, BIGINT, NUMERIC, TEXT, CITEXT, TEXT, BYTEA, BOOLEAN, BOOLEAN,
    TEXT, BOOLEAN, TEXT
);
