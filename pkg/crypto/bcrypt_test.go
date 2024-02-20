package crypto_test

import (
	"testing"

	"github.com/memnix/memnix-rest/pkg/crypto"
)

func TestBcryptCrypto_HashAndVerify(t *testing.T) {
	// Test 1: Correctly hashed password should verify successfully
	t.Run("Test 1 - Successful Hash and Verify", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password := "my_password"
		hashedPassword, _ := bc.Hash(password)
		match, err := bc.Verify(password, hashedPassword)
		if !match || err != nil {
			t.Errorf("Test 1: Verify with correct password failed")
		}
	})

	// Test 2: Incorrect password should not match
	t.Run("Test 2 - Incorrect Password Verification", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password := "my_password"
		hashedPassword, _ := bc.Hash(password)
		incorrectPassword := "incorrect_password"
		match, _ := bc.Verify(incorrectPassword, hashedPassword)
		if match {
			t.Errorf("Test 2: Verify with incorrect password should fail")
		}
	})

	// Test 3: Verify empty password
	t.Run("Test 3 - Empty Password Verification", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		hashedPassword, _ := bc.Hash("")
		match, err := bc.Verify("", hashedPassword)
		if !match || err != nil {
			t.Errorf("Test 3: Verify with empty password failed")
		}
	})

	// Test 4: Successful hash and verify with different cost
	t.Run("Test 4 - Hash and Verify with Different Cost", func(t *testing.T) {
		cost := 12
		bc := crypto.NewBcryptCrypto(cost)
		password := "my_password"
		hashedPassword, _ := bc.Hash(password)
		match, err := bc.Verify(password, hashedPassword)
		if !match || err != nil {
			t.Errorf("Test 4: Verify with correct password and different cost failed")
		}
	})

	// Test 5: Verify against empty hashed password
	t.Run("Test 5 - Verify Against Empty Hashed Password", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		incorrectPassword := "incorrect_password"
		emptyHashedPassword := []byte{}
		match, _ := bc.Verify(incorrectPassword, emptyHashedPassword)
		if match {
			t.Errorf("Test 5: Verify against empty hashed password should fail")
		}
	})

	// Test 6: Incorrect password with correct hashed password
	t.Run("Test 6 - Incorrect Password with Correct Hashed Password", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password := "my_password"
		hashedPassword, _ := bc.Hash(password)
		incorrectPassword := "incorrect_password"
		match, _ := bc.Verify(incorrectPassword, hashedPassword)
		if match {
			t.Errorf("Test 6: Incorrect password should not match the correct hashed password")
		}
	})

	// Test 7: Compare two hashed passwords
	t.Run("Test 7 - Compare Two Hashed Passwords", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password1 := "password1"
		password2 := "password2"
		_, _ = bc.Hash(password1)
		hashedPassword2, _ := bc.Hash(password2)
		match, _ := bc.Verify(password1, hashedPassword2)
		if match {
			t.Errorf("Test 7: Two different hashed passwords should not match")
		}
	})

	// Test 8: Successful hash and verify with empty password
	t.Run("Test 8 - Hash and Verify with Empty Password", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		emptyPassword := ""
		hashedPassword, _ := bc.Hash(emptyPassword)
		match, err := bc.Verify(emptyPassword, hashedPassword)
		if !match || err != nil {
			t.Errorf("Test 8: Verify with empty password should succeed")
		}
	})

	// Test 9: Verify against incorrect hashed password
	t.Run("Test 9 - Verify Against Incorrect Hashed Password", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password := "my_password"
		_, _ = bc.Hash(password)
		incorrectHashedPassword := []byte("incorrect_hashed_password")
		match, _ := bc.Verify(password, incorrectHashedPassword)
		if match {
			t.Errorf("Test 9: Verify against incorrect hashed password should fail")
		}
	})

	// Test 10: BcryptCrypto should fail when the password exceeds 72 bytes
	t.Run("Test 10 - Password Length Exceeds 72 Bytes", func(t *testing.T) {
		cost := 10
		bc := crypto.NewBcryptCrypto(cost)
		password := "ThisPasswordExceedsSeventyTwoBytesToFailTheTestPasswordIsLimitedToSeventyTwoCharacters"
		_, err := bc.Hash(password)
		if err == nil {
			t.Errorf("Test 10: Hashing a password exceeding 72 bytes should fail")
		}
	})
}

func runBenchmark(b *testing.B, cost int) {
	bc := crypto.NewBcryptCrypto(cost)
	password := "my_password"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bc.Hash(password)
	}
}

func BenchmarkAllCosts(b *testing.B) {
	// Define the bcrypt cost values to benchmark
	costs := []struct {
		Name string
		Cost int
	}{
		{"Cost4", 4},
		{"Cost10", 10},
		{"Cost12", 12},
		{"Cost14", 14},
	}

	// Run benchmarks in parallel
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, cost := range costs {
				b.Run(cost.Name, func(b *testing.B) {
					runBenchmark(b, cost.Cost)
				})
			}
		}
	})
}
