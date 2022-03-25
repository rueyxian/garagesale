package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// ================================================================================
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// ================================================================================
// Claims
type Claims struct {
	Roles []string
	jwt.StandardClaims
}

// ================================================================================
// NewClaims
func NewClaims(subject string, roles []string, now time.Time, expires time.Duration) Claims {
	return Claims{
		Roles: roles,
		StandardClaims{
			Subject:   subject,
			IssueAt:   now.Unix(),
			ExpiresAt: now.Add(expires).Unix(),
		},
	}
}
