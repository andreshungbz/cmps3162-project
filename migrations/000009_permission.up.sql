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
    ('hotel:read'), -- 1
    ('hotel:write'), -- 2
    ('department:read'), -- 3
    ('department:write'), -- 4
    ('room_type:read'), -- 5
    ('room_type:write'), -- 6
    ('room:read'), -- 7
    ('room:write'), -- 8
    ('employee:read'), -- 9
    ('employee:write'), -- 10
    ('housekeeping_task:read'), -- 11
    ('housekeeping_task:write'), -- 12
    ('maintenance_report:read'), -- 13
    ('maintenance_report:write'), -- 14
    ('guest:read'), -- 15
    ('guest:write'), -- 16
    ('reservation:read'), -- 17
    ('reservation:write'), -- 18
    ('registration:read'), -- 19
    ('registration:write'); -- 20

INSERT INTO employee_permission(employee_id, permission_id)
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (1, 5),
    (1, 6),
    (1, 7),
    (1, 8),
    (1, 9),
    (1, 10),

    (2, 7),
    (2, 8),
    (2, 15),
    (2, 16),
    (2, 17),
    (2, 18),
    (2, 19),
    (2, 20),

    (3, 7),
    (3, 8),
    (3, 11),
    (3, 12),
    (3, 13),
    (3, 14);
