package mock

import (
	"time"

	"github.com/victorprb/goapi/pkg/models"
)

var mockUser = &models.User{
	UUID:      "29c31a03-f111-4f6f-b49e-07b0117fbb42",
	FirstName: "Cloud",
	LastName:  "Strife",
	Email:     "cloud.strife@ffvii.com",
	CPF:       "12312312300",
	Created:   time.Now(),
}

var mockCredential = &models.Credentials{
	Login:    "cloud.strife@ffvii.com",
	Password: "my_app",
}

type UserModel struct{}

func (m *UserModel) Insert(u *models.User) error {
	switch u.Email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(c *models.Credentials) (string, error) {
	switch c.Login {
	case "cloud.strife@ffvii.com":
		return "29c31a03-f111-4f6f-b49e-07b0117fbb42", nil
	default:
		return "", models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(uuid string) (*models.User, error) {
	switch uuid {
	case "29c31a03-f111-4f6f-b49e-07b0117fbb42":
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
