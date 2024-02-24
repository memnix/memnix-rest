package domain_test

import (
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestPermission_IsValid(t *testing.T) {
	tests := []struct {
		permission domain.Permission
		expected   bool
	}{
		{permission: domain.PermissionNone, expected: true},
		{permission: domain.PermissionUser, expected: true},
		{permission: domain.PermissionVip, expected: true},
		{permission: domain.PermissionAdmin, expected: true},
		{permission: -1, expected: false},
		{permission: 4, expected: false},
	}

	for _, test := range tests {
		result := test.permission.IsValid()
		if result != test.expected {
			t.Errorf("Expected IsValid() to return %v for permission %d, but got %v", test.expected, test.permission, result)
		}
	}
}

func TestPermission_String(t *testing.T) {
	tests := []struct {
		expected   string
		permission domain.Permission
	}{
		{permission: domain.PermissionNone, expected: "none"},
		{permission: domain.PermissionUser, expected: "user"},
		{permission: domain.PermissionVip, expected: "vip"},
		{permission: domain.PermissionAdmin, expected: "admin"},
	}

	for _, test := range tests {
		result := test.permission.String()
		if result != test.expected {
			t.Errorf("Expected String() to return %s for permission %d, but got %s", test.expected, test.permission, result)
		}
	}
}
