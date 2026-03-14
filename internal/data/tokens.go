package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"time"

	"github.com/andreshungbz/cmps3162-project/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

// Token maps the token entity.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	PersonID  int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// generateToken generates a token from a SHA-546 hash of a randomly generated text string.
func generateToken(personID int64, ttl time.Duration, scope string) *Token {
	token := &Token{
		Plaintext: rand.Text(),
		PersonID:  personID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token
}

// ValidateTokenPlaintext validates the length of a token.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// TokenModel holds the database handler.
type TokenModel struct {
	DB *sql.DB
}

// New generates a token and inserts it into the database.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := generateToken(userID, ttl, scope)

	err := m.Insert(token)
	return token, err
}

// Insert creates a token record.
func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO token (hash, person_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.PersonID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForPerson removes all token records by a person's ID and scope.
func (m TokenModel) DeleteAllForPerson(scope string, personID int64) error {
	query := `
		DELETE FROM token
		WHERE scope = $1 AND person_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, personID)
	return err
}
