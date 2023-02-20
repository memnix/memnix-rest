package env

import "os"

// IEnv is an interface for environment variables
type IEnv interface {
	GetEnv(key string) string // GetEnv returns the value of the environment variable named by the key.
}

// Env is a struct for environment variables
type Env struct {
	env IEnv // env is an interface for environment variables
}

// NewMyEnv returns a new Env struct
func NewMyEnv(env IEnv) *Env {
	return &Env{env: env}
}

// GetEnv returns the value of the environment variable named by the key.
func (m *Env) GetEnv(key string) string {
	return m.env.GetEnv(key)
}

// FakeEnv is a struct for fake environment variables
type FakeEnv struct{}

// GetEnv returns the value of the environment variable named by the key.
// It returns predefined value for testing.
func (*FakeEnv) GetEnv(key string) string {
	fakeEnv := map[string]string{
		"APP_ENV":    "dev",    // development
		"SECRET_KEY": "secret", // secret key
	}
	// return predefined value for existing key
	if val, ok := fakeEnv[key]; ok {
		return val
	}
	// return empty string for non-existing key
	return ""
}

// OsEnv is a struct for os environment variables
type OsEnv struct{}

// GetEnv returns the value of the environment variable named by the key.
func (*OsEnv) GetEnv(key string) string {
	return os.Getenv(key) // return os environment variable
}
