-- migrations/token_table.up.sql
-- Creates the table for token activation.

CREATE TABLE IF NOT EXISTS token (
    hash BYTEA PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);
