package models

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateCPF       = errors.New("models: duplicate cpf")
)

// Credentials model
type Credentials struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

// Claims model
type Claims struct {
	Login string `json:"login,omitempty"`
	jwt.StandardClaims
}

// User model
type User struct {
	UUID           string    `json:"uuid,omitempty"`
	FirstName      string    `json:"firstname,omitempty"`
	LastName       string    `json:"lastname,omitempty"`
	CPF            string    `json:"cpf,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"password,omitempty"`
	Created        time.Time `json:"created,omitempty"`
	Active         bool      `json:"active,omitempty"`
}
