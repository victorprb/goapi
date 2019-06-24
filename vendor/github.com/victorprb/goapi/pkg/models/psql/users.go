package psql

import (
	"database/sql"
	"strings"

	"github.com/victorprb/goapi/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserModel db for app
type UserModel struct {
	DB *sql.DB
}

// Insert a new user
func (m *UserModel) Insert(u *models.User) error {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.HashedPassword), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (first_name, last_name, cpf, email, hashed_password, created)
	VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING uuid`

	err = m.DB.QueryRow(stmt, u.FirstName, u.LastName, u.CPF, u.Email, string(HashedPassword)).Scan(&u.UUID)
	if err != nil {
		if psqlErr, ok := err.(*pq.Error); ok {
			if psqlErr.Code.Name() == "unique_violation" {
				if strings.Contains(psqlErr.Constraint, "users_email_key") {
					return models.ErrDuplicateEmail
				}
				return models.ErrDuplicateCPF
			}
		}
	}
	return err
}

// Authenticate check user credentials
func (m *UserModel) Authenticate(c *models.Credentials) (string, error) {
	var uuid string
	var hashedPassword []byte
	stmt := "SELECT uuid, hashed_password FROM users WHERE email = $1 AND active = true"
	row := m.DB.QueryRow(stmt, c.Login)
	err := row.Scan(&uuid, &hashedPassword)
	if err == sql.ErrNoRows {
		return "", models.ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(c.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", models.ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	return uuid, nil
}

// Get user info
func (m *UserModel) Get(uuid string) (*models.User, error) {
	u := &models.User{}

	stmt := "SELECT uuid, first_name, last_name, cpf, email, created, active FROM users WHERE uuid = $1 AND active = true"
	err := m.DB.QueryRow(stmt, uuid).Scan(&u.UUID, &u.FirstName, &u.LastName, &u.CPF, &u.Email, &u.Created, &u.Active)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return u, nil
}
