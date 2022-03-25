package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// ================================================================================

var (
	ErrAuthenticationFailure = errors.New("Authentication failed")
)

// ================================================================================
// Create
func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "hashing user password")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: passHash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users
		(user_id, name, email, roles, password_hash, date_created, date_updated)
		VALUES($1, $2, $3, $4, $5, $6, $7)`

	if _, err := db.ExecContext(
		ctx, q,
		u.ID, u.Name, u.Email, u.Roles, u.PasswordHash, u.DateCreated, u.DateUpdated,
	); err != nil {
		return nil, errors.Wrapf(err, "inserting user: %v", nu)
	}

	return &u, nil

}

// ================================================================================
// Authenticate
func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, email, password string) (authClaims, error) {

	const q = `SELECT * FROM users WHERE email = $1`

	var u User
	if err := db.GetContext(ctx, u, q); err != nil {
		if err == sql.ErrNoRows {
			return Claims{}, ErrAuthenticationFailure
		}
		return Claims{}, errors.Wrap(err, "selecting single user")
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return Claims{}, ErrAuthenticationFailure
	}

	c := NewClaims(u.ID, u.Roles, now, time.Hour)
	return c, nil
}

// ================================================================================
