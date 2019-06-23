package psql

import (
	"database/sql"
	"strings"

	"github.com/victorprb/goapi/pkg/models"

	"github.com/lib/pq"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(u *models.User) error {
	stmt := `INSERT INTO users (first_name, last_name, cpf, email, created)
	VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING uuid`

	err := m.DB.QueryRow(stmt, u.FirstName, u.LastName, u.CPF, u.Email).Scan(&u.UUID)
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
