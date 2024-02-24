package config_test

import (
	"testing"

	"github.com/memnix/memnix-rest/cmd/v2/config"
)

func TestGetConfigPath(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult string
		isDevelopment  bool
	}{
		{
			name:           "DevelopmentEnvironment",
			isDevelopment:  true,
			expectedResult: "./cmd/v2/config/config-local",
		},
		{
			name:           "ProductionEnvironment",
			isDevelopment:  false,
			expectedResult: "./cmd/v2/config/config-prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetConfigPath(tt.isDevelopment)
			if result != tt.expectedResult {
				t.Errorf("GetConfigPath() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		expectedResult bool
	}{
		{
			name:           "ProductionEnvironment",
			env:            "prod",
			expectedResult: true,
		},
		{
			name:           "DevelopmentEnvironment",
			env:            "dev",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.env)
			result := config.IsProduction()
			if result != tt.expectedResult {
				t.Errorf("IsProduction() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		expectedResult bool
	}{
		{
			name:           "ProductionEnvironment",
			env:            "prod",
			expectedResult: false,
		},
		{
			name:           "DevelopmentEnvironment",
			env:            "dev",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.env)
			result := config.IsDevelopment()
			if result != tt.expectedResult {
				t.Errorf("IsDevelopment() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGetCallbackURL(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		expectedResult string
	}{
		{
			name:           "CallbackURLSet",
			env:            "https://example.com/callback",
			expectedResult: "https://example.com/callback",
		},
		{
			name:           "CallbackURLNotSet",
			env:            "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CALLBACK_URL", tt.env)
			result := config.GetCallbackURL()
			if result != tt.expectedResult {
				t.Errorf("GetCallbackURL() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
