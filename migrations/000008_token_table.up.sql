-- migrations/token_table.up.sql
-- Creates the table for token activation.

CREATE TABLE IF NOT EXISTS token (
    hash BYTEA PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES person(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);

INSERT INTO token (hash, person_id, expiry, scope)
VALUES
    (digest('NTH5BYOPVS6WJ2XOK27W7Q6ZUH', 'sha256'), 1, '2126-05-15', 'authentication'),
    (digest('UZPHDJXKVJORSLLQD7AH3QDCU3', 'sha256'), 2, '2126-05-15', 'authentication'),
    (digest('UKAKGCKBNSDEJPYBN7S6LFZI3W', 'sha256'), 3, '2126-05-15', 'authentication');
