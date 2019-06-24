package psql

import (
	"reflect"
	"testing"
	"time"

	"github.com/victorprb/goapi/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("psql: skipping integration test")
	}

	tests := []struct {
		name string
		userUUID string
		wantUser *models.User
		wantError error
	}{
		{
			name: "Valid UUID",
			userUUID: "29c31a03-f111-4f6f-b49e-07b0117fbb42",
			wantUser: &models.User{
				UUID: "29c31a03-f111-4f6f-b49e-07b0117fbb42",
				FirstName: "Cloud",
				LastName: "Strife",
				CPF: "12312312300",
				Email: "cloud.strife@ffvii.com",
				Created: time.Date(2019, 06, 20, 15, 0, 0, 0, time.UTC),
				Active: true,
			},
			wantError: nil,
		},
		{
			name: "Non-existent UUID",
			userUUID: "3aa9289e-d80f-4a67-bc77-2b68271fb98b",
			wantUser: nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize a connection pool to our test database
			db, teardown := newTestDB(t)
			defer teardown()

			// Create a new instance of the UserModel
			m := UserModel{db}

			user, err := m.Get(tt.userUUID)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %#v; got %#v", tt.wantUser, user)
			}
		})
	}
}