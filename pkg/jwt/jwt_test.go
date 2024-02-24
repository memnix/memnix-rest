package jwt_test

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/pkg/crypto"
	myjwt "github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

func TestGetJwtInstance(t *testing.T) {
	// Create a wait group to ensure the function is called only once
	var wg sync.WaitGroup
	wg.Add(2)

	// Call the function in two goroutines
	go func() {
		defer wg.Done()
		instance1 := myjwt.GetJwtInstance()
		if instance1 == nil {
			t.Error("GetJwtInstance() returned nil")
		}
	}()

	go func() {
		defer wg.Done()
		instance2 := myjwt.GetJwtInstance()
		if instance2 == nil {
			t.Error("GetJwtInstance() returned nil")
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}

func TestNewJWTInstance(t *testing.T) {
	headerLen := 10
	expirationTime := 24
	publicKey := make(ed25519.PublicKey, 32)
	privateKey := make(ed25519.PrivateKey, 64)

	instance := myjwt.NewJWTInstance(headerLen, expirationTime, publicKey, privateKey)

	if instance.HeaderLen != headerLen {
		t.Errorf("Expected HeaderLen to be %d, but got %d", headerLen, instance.HeaderLen)
	}

	if !bytes.Equal(instance.PublicKey, publicKey) {
		t.Error("Expected PublicKey to be equal")
	}

	if !bytes.Equal(instance.PrivateKey, privateKey) {
		t.Error("Expected PrivateKey to be equal")
	}

	if instance.SigningMethod != jwt.SigningMethodEdDSA {
		t.Errorf("Expected signingMethod to be %s, but got %s", jwt.SigningMethodEdDSA, instance.SigningMethod)
	}

	if instance.ExpirationTimeInHours != expirationTime {
		t.Errorf("Expected ExpirationTimeInHours to be %d, but got %d", expirationTime, instance.ExpirationTimeInHours)
	}
}

func TestSetJwt(t *testing.T) {
	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(10, 24, make(ed25519.PublicKey, 32), make(ed25519.PrivateKey, 64))

	// Create a new InstanceSingleton
	singleton := &myjwt.InstanceSingleton{}

	// Set the JWT instance
	singleton.SetJwt(instance)

	// Verify that the JWT instance is set correctly
	if !bytes.Equal(singleton.GetJwt().PublicKey, instance.PublicKey) {
		t.Error("SetJwt() did not set the JWT instance correctly")
	}

	if !bytes.Equal(singleton.GetJwt().PrivateKey, instance.PrivateKey) {
		t.Error("SetJwt() did not set the JWT instance correctly")
	}

	if singleton.GetJwt().HeaderLen != instance.HeaderLen {
		t.Error("SetJwt() did not set the JWT instance correctly")
	}

	if singleton.GetJwt().ExpirationTimeInHours != instance.ExpirationTimeInHours {
		t.Error("SetJwt() did not set the JWT instance correctly")
	}

	if singleton.GetJwt().SigningMethod != instance.SigningMethod {
		t.Error("SetJwt() did not set the JWT instance correctly")
	}
}

func TestGetJwt(t *testing.T) {
	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(10, 24, make(ed25519.PublicKey, 32), make(ed25519.PrivateKey, 64))

	// Create a new InstanceSingleton
	singleton := &myjwt.InstanceSingleton{}

	// Set the JWT instance
	singleton.SetJwt(instance)

	// Verify that the JWT instance is returned correctly
	if !bytes.Equal(singleton.GetJwt().PublicKey, instance.PublicKey) {
		t.Error("GetJwt() did not return the JWT instance correctly")
	}

	if !bytes.Equal(singleton.GetJwt().PrivateKey, instance.PrivateKey) {
		t.Error("GetJwt() did not return the JWT instance correctly")
	}

	if singleton.GetJwt().HeaderLen != instance.HeaderLen {
		t.Error("GetJwt() did not return the JWT instance correctly")
	}

	if singleton.GetJwt().ExpirationTimeInHours != instance.ExpirationTimeInHours {
		t.Error("GetJwt() did not return the JWT instance correctly")
	}

	if singleton.GetJwt().SigningMethod != instance.SigningMethod {
		t.Error("GetJwt() did not return the JWT instance correctly")
	}
}

func ParseEd25519Key() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	publicKey, privateKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func TestGenerateToken(t *testing.T) {
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)

	// Create a mock context
	ctx := context.TODO()

	// Set up test cases
	testCases := []struct {
		name   string
		userID uint
	}{
		{
			name:   "ValidUserID",
			userID: 1,
		},
		{
			name:   "InvalidUserID",
			userID: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := instance.GenerateToken(ctx, tc.userID)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if token == "" {
				t.Error("Expected token to be non-empty")
			}

			jwtToken, err := instance.GetToken(ctx, token)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check if the token is valid.
			userID, err := instance.VerifyToken(jwtToken)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if userID != tc.userID {
				t.Errorf("Expected userID to be %d, but got %d", tc.userID, userID)
			}

			// Get the expiration time from the token
			expirationTime := instance.GetExpirationTime(jwtToken)

			if expirationTime == 0 {
				t.Error("Expected expirationTime to be non-zero")
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)
	token, err := instance.GenerateToken(context.TODO(), 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	testCases := []struct {
		name        string
		token       string
		expected    string
		expectedLen int
	}{
		{
			name:     "ValidToken",
			token:    "Bearer " + token,
			expected: token,
		},
		{
			name:     "InvalidToken",
			token:    "Bearer",
			expected: "",
		},
		// Add more test cases here if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := instance.ExtractToken(tc.token)
			if result != tc.expected {
				t.Errorf("Expected token to be %s, but got %s", tc.expected, result)
			}
		})
	}
}

func GenerateInvalidToken() string {
	// Create a new JWT instance
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		return ""
	}
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)

	// Create the Claims for the token
	claims := jwt.NewWithClaims(instance.SigningMethod, jwt.RegisteredClaims{
		Issuer:    utils.ConvertUIntToStr(42),         // Issuer is the user id
		ExpiresAt: instance.CalculateExpirationTime(), // ExpiresAt is the expiration time
	})

	// Sign and get the complete encoded token as a string using the secret
	token, err := claims.SignedString(instance.PrivateKey)
	if err != nil {
		return ""
	}
	return token
}

func GenerateTokenWithString(id string) string {
	// Create a new JWT instance
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		return ""
	}
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)

	// Create the Claims for the token
	claims := jwt.NewWithClaims(instance.SigningMethod, jwt.RegisteredClaims{
		Issuer:    id,
		ExpiresAt: instance.CalculateExpirationTime(), // ExpiresAt is the expiration time
	})

	// Sign and get the complete encoded token as a string using the secret
	token, err := claims.SignedString(instance.PrivateKey)
	if err != nil {
		return ""
	}
	return token
}

func TestGetConnectedUserID(t *testing.T) {
	// Create a mock context
	ctx := context.TODO()

	// Create a new JWT instance
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)

	token, err := instance.GenerateToken(ctx, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	zeroToken, err := instance.GenerateToken(ctx, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Set up test cases
	testCases := []struct {
		expectedErr error
		name        string
		tokenHeader string
		expectedID  uint
	}{
		{
			name:        "ValidToken",
			tokenHeader: "Bearer " + token,
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name:        "InvalidToken",
			tokenHeader: "Bearer " + GenerateInvalidToken(),
			expectedID:  0,
			expectedErr: errors.New("invalid token"),
		},
		{
			name:        "EmptyToken",
			tokenHeader: "",
			expectedID:  0,
			expectedErr: errors.New("empty token"),
		},
		{
			name:        "InvalidTokenHeader",
			tokenHeader: "Bearer",
			expectedID:  0,
			expectedErr: errors.New("empty token"),
		},
		{
			name:        "ZeroToken",
			tokenHeader: "Bearer " + zeroToken,
			expectedID:  0,
			expectedErr: errors.New("invalid token"),
		},
		{
			name:        "StringToken",
			tokenHeader: "Bearer " + GenerateTokenWithString("42"),
			expectedID:  0,
			expectedErr: errors.New("invalid token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := instance.GetConnectedUserID(ctx, tc.tokenHeader)

			if userID != tc.expectedID {
				t.Errorf("Expected userID to be %d, but got %d", tc.expectedID, userID)
			}

			if (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr == nil) || (err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Expected error to be %v, but got %v", tc.expectedErr, err)
			}
		})
	}
}

func FuzzGenerateToken(f *testing.F) {
	privateKey, publicKey, err := ParseEd25519Key()
	if err != nil {
		f.Fatalf("Unexpected error: %v", err)
	}

	// Create a new JWT instance
	instance := myjwt.NewJWTInstance(2, 24, publicKey, privateKey)

	// Create a mock context
	ctx := context.TODO()

	f.Add(uint(1)) // Add a value to the fuzzing pool

	f.Fuzz(func(t *testing.T, userID uint) {
		token, err := instance.GenerateToken(ctx, userID)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if token == "" {
			t.Error("Expected token to be non-empty")
		}

		jwtToken, err := instance.GetToken(ctx, token)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Check if the token is valid.
		uID, err := instance.VerifyToken(jwtToken)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if userID != uID {
			t.Errorf("Expected userID to be %d, but got %d", userID, uID)
		}
	})
}
