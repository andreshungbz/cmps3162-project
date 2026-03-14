-- migrations/permission.up.sql
-- Creates the table for permissions.

CREATE TABLE IF NOT EXISTS permission (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS employee_permission (
    employee_id BIGINT NOT NULL REFERENCES employee ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permission ON DELETE CASCADE,
    PRIMARY KEY (employee_id, permission_id)
);

INSERT INTO permission(code)
VALUES
    ('housekeeping_task:read'),
    ('housekeeping_task:write'),
    ('maintenance_report:read'),
    ('maintenance_report:write');

INSERT INTO employee_permission(employee_id, permission_id)
VALUES
    (1, 1),
    (1, 3),
    (3, 1),
    (3, 2),
    (3, 3),
    (3, 4);
