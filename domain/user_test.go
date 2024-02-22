package domain_test

import (
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestUser_ToPublicUser(t *testing.T) {
	user := &domain.User{
		Username:   "testuser",
		Email:      "test@example.com",
		Avatar:     "https://example.com/avatar.png",
		Permission: domain.PermissionUser,
	}

	user.ID = 1

	expected := domain.PublicUser{
		Username:   "testuser",
		Email:      "test@example.com",
		Avatar:     "https://example.com/avatar.png",
		ID:         1,
		Permission: domain.PermissionUser,
	}

	publicUser := user.ToPublicUser()

	if publicUser != expected {
		t.Errorf("Expected public user to be %v, but got %v", expected, publicUser)
	}
}

func TestUser_HasPermission(t *testing.T) {
	useCases := []struct {
		user     *domain.User
		name     string
		perm     domain.Permission
		expected bool
	}{
		{
			name: "UserHasPermission",
			user: &domain.User{
				Permission: domain.PermissionUser,
			},
			perm:     domain.PermissionUser,
			expected: true,
		},
		{
			name: "UserDoesNotHavePermission",
			user: &domain.User{
				Permission: domain.PermissionUser,
			},
			perm:     domain.PermissionAdmin,
			expected: false,
		},
		{
			name: "AdminHasPermission",
			user: &domain.User{
				Permission: domain.PermissionAdmin,
			},
			perm:     domain.PermissionUser,
			expected: true,
		},
		{
			name: "AdminHasAdminPermission",
			user: &domain.User{
				Permission: domain.PermissionAdmin,
			},
			perm:     domain.PermissionAdmin,
			expected: true,
		},
	}

	for _, uc := range useCases {
		t.Run(uc.name, func(t *testing.T) {
			result := uc.user.HasPermission(uc.perm)
			if result != uc.expected {
				t.Errorf("Expected result to be %v, but got %v", uc.expected, result)
			}
		})
	}
}

func TestRegister_Validate(t *testing.T) {
	testCases := []struct {
		register  *domain.Register
		name      string
		shouldErr bool
	}{
		{
			name: "Valid Register",
			register: &domain.Register{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			shouldErr: false,
		},
		{
			name: "No Username Register",
			register: &domain.Register{
				Username: "",
				Email:    "test@example.com",
				Password: "password",
			},
			shouldErr: true,
		},
		{
			name: "No Email Register",
			register: &domain.Register{
				Username: "testuser",
				Email:    "",
				Password: "password",
			},
			shouldErr: true,
		},
		{
			name: "No Password Register",
			register: &domain.Register{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.register.Validate()
			if tc.shouldErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}
		})
	}
}

func TestRegister_ToUser(t *testing.T) {
	register := &domain.Register{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	expected := domain.User{
		Username:   "testuser",
		Email:      "test@example.com",
		Password:   "password",
		Permission: domain.PermissionUser,
	}

	user := register.ToUser()

	if user.Username != expected.Username {
		t.Errorf("Expected user username to be %s, but got %s", expected.Username, user.Username)
	}

	if user.Email != expected.Email {
		t.Errorf("Expected user email to be %s, but got %s", expected.Email, user.Email)
	}

	if user.Password != expected.Password {
		t.Errorf("Expected user password to be %s, but got %s", expected.Password, user.Password)
	}

	if user.Permission != expected.Permission {
		t.Errorf("Expected user permission to be %v, but got %v", expected.Permission, user.Permission)
	}
}

func TestUser_TableName(t *testing.T) {
	user := &domain.User{}
	expected := "users"
	result := user.TableName()
	if result != expected {
		t.Errorf("Expected table name to be %s, but got %s", expected, result)
	}
}

func TestUser_Validate(t *testing.T) {
	useCases := []struct {
		user      *domain.User
		name      string
		shouldErr bool
	}{
		{
			name: "Valid User",
			user: &domain.User{
				Username: "testuser",
				Email:    "toto@toto.com",
				Password: "password",
			},
			shouldErr: false,
		},
		{
			name: "No Username User",
			user: &domain.User{
				Username: "",
				Email:    "toto@toto.com",
				Password: "password",
			},
			shouldErr: true,
		},
		{
			name: "No Email User",
			user: &domain.User{
				Username: "testuser",
				Email:    "",
				Password: "password",
			},
			shouldErr: true,
		},
		{
			name: "No Password User",
			user: &domain.User{
				Username: "testuser",
				Email:    "toto@toto.com",
				Password: "",
			},
			shouldErr: true,
		},
	}

	for _, uc := range useCases {
		t.Run(uc.name, func(t *testing.T) {
			err := uc.user.Validate()
			if uc.shouldErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !uc.shouldErr && err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}
		})
	}
}
