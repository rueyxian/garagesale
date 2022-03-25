package user

import (
	"time"

	"github.com/lib/pq"
)

// User
type User struct {
	ID           string         `json:"id" db:"user_id"`
	Name         string         `json:"name" db:"name"`
	Email        string         `json:"email" db:"email"`
	Roles        pq.StringArray `json:"roles" db:"roles"`
	PasswordHash []byte         `json:"-" db:"password_hash"`
	DateCreated  time.Time      `json:"date_created" db:"date_created"`
	DateUpdated  time.Time      `json:"date_updated" db:"date_updated"`
}

// NewUser
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}
