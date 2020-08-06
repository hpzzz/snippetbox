package mysql

import (
	"reflect"
	"testing"
	"time"

	"karolharasim.com/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	//Skip the test if the '-short' flag is provided when running the test.
	if testing.Short() {
		t.Skip("mysql: skipping intergation test")
	}
	
	tests := []struct {
		name string
		userID int
		wantUser *models.User
		wantError error
	} {
		{
			name: "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID: 1,
				Name: "Kari Hari",
				Email: "Kari@example.com",
				Created: time.Date(2020, 07, 23, 18, 12, 0, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name: "Zero ID",
			userID: 0,
			wantUser: nil,
			wantError: models.ErrNoRecord,
		},
		{
			name: "Non-existent ID",
			userID: 2,
			wantUser: nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Initialize a connection pool to our test database and defer a call to
			// the teardown function so it is always run imediately
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}
			
			user, err := m.Get(tt.userID)

			if err != tt.wantError {
				t.Errorf("want: %v, got: %v", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want: %v, got: %v", tt.wantUser, user)
			}
		})
	}
}