package random_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/memnix/memnix-rest/pkg/random"
)

func TestGenerator_GenerateSecretCode(t *testing.T) {
	generator := random.GetRandomGeneratorInstance()

	testCases := []struct {
		expectedErr error
		name        string
		n           int
		expected    int
		hasErr      bool
	}{
		{
			name:        "GenerateSecretCodeWithLength10",
			n:           10,
			expected:    10,
			expectedErr: nil,
			hasErr:      false,
		},
		{
			name:        "GenerateSecretCodeWithLength20",
			n:           20,
			expected:    20,
			expectedErr: nil,
			hasErr:      false,
		},
		{
			name:        "GenerateSecretCodeWithLength0",
			n:           0,
			expected:    0,
			expectedErr: nil,
			hasErr:      false,
		},
		{
			name:        "GenerateSecretCodeWithLengthNegative",
			n:           -1,
			expected:    0,
			expectedErr: errors.New("invalid length"),
			hasErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code, err := generator.GenerateSecretCode(tc.n)
			if tc.hasErr && err == nil {
				t.Fatalf("Expected error, but got nil")
			}

			if !tc.hasErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tc.hasErr && err != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("Expected error to be %v, but got %v", tc.expectedErr, err)
			}

			if len(code) != tc.expected {
				t.Errorf("Expected code length to be %d, but got %d", tc.expected, len(code))
			}

			if !isValidCode(code) {
				t.Errorf("Generated code is not valid: %s", code)
			}
		})
	}
}

func isValidCode(code string) bool {
	// Define the valid characters for the secret code
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Check if each character in the code is a valid character
	for _, c := range code {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}

	return true
}

func FuzzGenerateSecretCode(f *testing.F) {
	generator := random.GetRandomGeneratorInstance()
	f.Add(uint(1000))
	f.Fuzz(func(t *testing.T, size uint) {
		code, err := generator.GenerateSecretCode(int(size))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(code) != int(size) {
			t.Errorf("Expected code length to be %d, but got %d", size, len(code))
		}

		if !isValidCode(code) {
			t.Errorf("Generated code is not valid: %s", code)
		}
	})
}
