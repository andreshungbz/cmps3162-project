-- migrations/employee_functions.up.sql
-- Creates functions for employees (subtype of person).

-- ====================================================================================
-- CREATE FUNCTION fn_create_employee returns the person id and created_at for a newly
-- created employee. Handles conditional insertion into role-specific tables.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_create_employee(
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT,
    -- employee attributes
    p_hotel_id BIGINT,
    p_department TEXT,
    p_manager_id BIGINT,
    p_salary NUMERIC,
    p_ssn TEXT,
    p_work_email CITEXT,
    p_work_phone TEXT,
    p_password_hash BYTEA,
    p_employed BOOLEAN,
    p_activated BOOLEAN,
    -- role-specific attributes
    p_role TEXT,  -- "operations_manager", "front_desk", "housekeeper"
    p_hotel_owner BOOLEAN DEFAULT NULL,
    p_shift TEXT DEFAULT NULL
)
RETURNS TABLE (
    id BIGINT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE
)
AS $$
DECLARE
    v_employee_id BIGINT;
BEGIN
    -- insert person entry
    INSERT INTO person (name, gender, street, city, country)
    VALUES (p_name, p_gender, p_street, p_city, p_country)
    RETURNING person.id INTO v_employee_id;

    -- insert employee entry
    INSERT INTO employee (id, hotel_id, department, manager_id, salary, ssn, work_email, work_phone, password_hash, employed, activated)
    VALUES (v_employee_id, p_hotel_id, p_department, p_manager_id, p_salary, p_ssn, p_work_email, p_work_phone, p_password_hash, p_employed, p_activated);

    -- conditional insertion into role-specific tables
    IF p_role = 'operations_manager' THEN
        INSERT INTO operations_manager (id, hotel_owner)
        VALUES (v_employee_id, p_hotel_owner);
    ELSIF p_role = 'front_desk' THEN
        INSERT INTO front_desk (id, shift)
        VALUES (v_employee_id, p_shift::shift_type);
    ELSIF p_role = 'housekeeper' THEN
        INSERT INTO housekeeper (id, shift)
        VALUES (v_employee_id, p_shift::shift_type);
    END IF;

    RETURN QUERY
    SELECT
        p.id,
        p.created_at,
        p.modified_at
    FROM person p
    WHERE p.id = v_employee_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- READ FUNCTION fn_get_employee_by_email returns employee and person data by work email.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_employee_by_email(
    p_work_email CITEXT
)
RETURNS TABLE (
    -- person attributes
    name TEXT,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE,
    -- employee attributes
    id BIGINT,
    hotel_id BIGINT,
    department TEXT,
    manager_id BIGINT,
    salary NUMERIC,
    ssn TEXT,
    work_email CITEXT,
    work_phone TEXT,
    password_hash BYTEA,
    employed BOOLEAN,
    activated BOOLEAN,
    -- role-specific attributes
    role TEXT,
    hotel_owner BOOLEAN,
    shift shift_type
)
AS $$
DECLARE
    v_role TEXT;
    v_employee_id BIGINT;
BEGIN
    -- find employee id from email
    SELECT e.id
    INTO v_employee_id
    FROM employee e
    WHERE e.work_email = p_work_email;

    IF v_employee_id IS NULL THEN
        RAISE EXCEPTION '[employee-not-found] Employee with email % does not exist', p_work_email;
    END IF;

    -- determine role
    IF EXISTS (SELECT 1 FROM operations_manager om WHERE om.id = v_employee_id) THEN
        v_role := 'operations_manager';
    ELSIF EXISTS (SELECT 1 FROM front_desk fd WHERE fd.id = v_employee_id) THEN
        v_role := 'front_desk';
    ELSIF EXISTS (SELECT 1 FROM housekeeper hk WHERE hk.id = v_employee_id) THEN
        v_role := 'housekeeper';
    ELSE
        RAISE EXCEPTION '[employee-not-found] Employee with email % has no valid role', p_work_email;
    END IF;

    RETURN QUERY
    SELECT
        -- person attributes
        p.name,
        p.gender,
        p.street,
        p.city,
        p.country,
        p.created_at,
        p.modified_at,
        -- employee attributes
        e.id,
        e.hotel_id,
        e.department,
        e.manager_id,
        e.salary,
        e.ssn,
        e.work_email,
        e.work_phone,
        e.password_hash,
        e.employed,
        e.activated,
        -- role-specific attributes
        v_role,
        om.hotel_owner,
        COALESCE(fd.shift, hk.shift) AS shift
    FROM employee e
    JOIN person p ON p.id = e.id
    LEFT JOIN operations_manager om ON om.id = e.id
    LEFT JOIN front_desk fd ON fd.id = e.id
    LEFT JOIN housekeeper hk ON hk.id = e.id
    WHERE e.id = v_employee_id;
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- READ FUNCTION fn_get_employee_for_token returns employee and person data for a
-- valid token hash and scope.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_employee_for_token(
    p_token_hash BYTEA,
    p_scope TEXT
)
RETURNS TABLE (
    -- person attributes
    name TEXT,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE,
    -- employee attributes
    id BIGINT,
    hotel_id BIGINT,
    department TEXT,
    manager_id BIGINT,
    salary NUMERIC,
    ssn TEXT,
    work_email CITEXT,
    work_phone TEXT,
    password_hash BYTEA,
    employed BOOLEAN,
    activated BOOLEAN,
    -- role-specific attributes
    role TEXT,
    hotel_owner BOOLEAN,
    shift shift_type
)
AS $$
BEGIN
    RETURN QUERY
    SELECT
        -- person attributes
        p.name,
        p.gender,
        p.street,
        p.city,
        p.country,
        p.created_at,
        p.modified_at,
        -- employee attributes
        e.id,
        e.hotel_id,
        e.department,
        e.manager_id,
        e.salary,
        e.ssn,
        e.work_email,
        e.work_phone,
        e.password_hash,
        e.employed,
        e.activated,
        -- role attributes
        CASE
            WHEN om.id IS NOT NULL THEN 'operations_manager'
            WHEN fd.id IS NOT NULL THEN 'front_desk'
            WHEN hk.id IS NOT NULL THEN 'housekeeper'
        END AS role,
        om.hotel_owner,
        COALESCE(fd.shift, hk.shift) AS shift
    FROM token t
    JOIN employee e ON e.id = t.person_id
    JOIN person p ON p.id = e.id
    LEFT JOIN operations_manager om ON om.id = e.id
    LEFT JOIN front_desk fd ON fd.id = e.id
    LEFT JOIN housekeeper hk ON hk.id = e.id
    WHERE
        t.hash = p_token_hash
        AND t.scope = p_scope
        AND t.expiry > NOW();
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- READ FUNCTION fn_get_employees returns employee records with optional role filter.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_get_employees(
    p_role TEXT DEFAULT ''
)
RETURNS TABLE (
    total_records BIGINT,
    -- person attributes
    name TEXT,
    gender TEXT,
    street TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMP(0) WITH TIME ZONE,
    modified_at TIMESTAMP(0) WITH TIME ZONE,
    -- employee attributes
    id BIGINT,
    hotel_id BIGINT,
    department TEXT,
    manager_id BIGINT,
    salary NUMERIC,
    ssn TEXT,
    work_email CITEXT,
    work_phone TEXT,
    password_hash BYTEA,
    employed BOOLEAN,
    activated BOOLEAN,
    -- role attributes
    role TEXT,
    hotel_owner BOOLEAN,
    shift shift_type
)
AS $$
BEGIN
    RETURN QUERY
    SELECT
        count(*) OVER(),
        -- person attributes
        p.name,
        p.gender,
        p.street,
        p.city,
        p.country,
        p.created_at,
        p.modified_at,
        -- employee attributes
        e.id,
        e.hotel_id,
        e.department,
        e.manager_id,
        e.salary,
        e.ssn,
        e.work_email,
        e.work_phone,
        e.password_hash,
        e.employed,
        e.activated,
        -- role-specific attributes
        CASE
            WHEN om.id IS NOT NULL THEN 'operations_manager'
            WHEN fd.id IS NOT NULL THEN 'front_desk'
            WHEN hk.id IS NOT NULL THEN 'housekeeper'
        END AS role,
        om.hotel_owner,
        COALESCE(fd.shift, hk.shift) AS shift
    FROM employee e
    JOIN person p ON p.id = e.id
    LEFT JOIN operations_manager om ON om.id = e.id
    LEFT JOIN front_desk fd ON fd.id = e.id
    LEFT JOIN housekeeper hk ON hk.id = e.id
    WHERE (
        p_role = '' OR
        (om.id IS NOT NULL AND p_role = 'operations_manager') OR
        (fd.id IS NOT NULL AND p_role = 'front_desk') OR
        (hk.id IS NOT NULL AND p_role = 'housekeeper')
    );
END;
$$ LANGUAGE plpgsql;

-- ====================================================================================
-- UPDATE FUNCTION fn_update_employee updates person, employee, and optional
-- role-specific fields.
-- ====================================================================================

CREATE OR REPLACE FUNCTION fn_update_employee(
    -- person attributes
    p_name TEXT,
    p_gender TEXT,
    p_street TEXT,
    p_city TEXT,
    p_country TEXT,
    -- employee attributes
    p_id BIGINT,
    p_hotel_id BIGINT,
    p_department TEXT,
    p_manager_id BIGINT,
    p_salary NUMERIC,
    p_ssn TEXT,
    p_work_email CITEXT,
    p_work_phone TEXT,
    p_password_hash BYTEA,
    p_employed BOOLEAN,
    p_activated BOOLEAN,
    -- role-specific attributes
    p_role TEXT,
    p_hotel_owner BOOLEAN DEFAULT NULL,
    p_shift TEXT DEFAULT NULL
)
RETURNS VOID
AS $$
BEGIN
    -- update person
    UPDATE person
    SET name = p_name,
        gender = p_gender,
        street = p_street,
        city = p_city,
        country = p_country
    WHERE id = p_id;

    -- update employee
    UPDATE employee
    SET hotel_id = p_hotel_id,
        department = p_department,
        manager_id = p_manager_id,
        salary = p_salary,
        ssn = p_ssn,
        work_email = p_work_email,
        work_phone = p_work_phone,
        employed = p_employed,
        activated = p_activated
    WHERE id = p_id;

    -- update role-specific tables
    IF p_role = 'operations_manager' THEN
        UPDATE operations_manager
        SET hotel_owner = p_hotel_owner
        WHERE id = p_id;
    ELSIF p_role = 'front_desk' THEN
        UPDATE front_desk
        SET shift = p_shift::shift_type
        WHERE id = p_id;
    ELSIF p_role = 'housekeeper' THEN
        UPDATE housekeeper
        SET shift = p_shift::shift_type
        WHERE id = p_id;
    END IF;
END;
$$ LANGUAGE plpgsql;
