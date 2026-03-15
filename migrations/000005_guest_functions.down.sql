-- migrations/guest_functions.down.sql
-- Drops all functions for guests in reverse order of their creation.

DROP FUNCTION IF EXISTS fn_delete_guest(TEXT);

DROP FUNCTION IF EXISTS fn_update_guest(
    TEXT, TEXT, TEXT, TEXT, TEXT,
    TEXT, CITEXT, TEXT
);

DROP FUNCTION IF EXISTS fn_get_guests(TEXT, TEXT);

DROP FUNCTION IF EXISTS fn_get_guest_by_passport(TEXT);

DROP FUNCTION IF EXISTS fn_create_guest(
    TEXT, TEXT, TEXT, TEXT, TEXT,
    TEXT, CITEXT, TEXT
);
