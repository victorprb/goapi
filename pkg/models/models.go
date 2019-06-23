package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord       = errors.New("models: no matching record found")
	ErrDuplicateEmail = errors.New("models: duplicate email")
	ErrDuplicateCPF   = errors.New("models: duplicate cpf")
)

// User model
type User struct {
	UUID      string    `json:"uuid"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created"`
	Active    bool      `json:"active"`
}
