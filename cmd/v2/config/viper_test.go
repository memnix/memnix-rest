package config_test

import (
	"log/slog"
	"testing"

	"github.com/memnix/memnix-rest/cmd/v2/config"
	"github.com/stretchr/testify/assert"
)

func TestLogConfig_GetSlogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected slog.Level
	}{
		{
			name:     "DebugLevel",
			level:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "InfoLevel",
			level:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "WarnLevel",
			level:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "ErrorLevel",
			level:    "error",
			expected: slog.LevelError,
		},
		{
			name:     "DefaultLevel",
			level:    "unknown",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logConfig := config.LogConfig{
				Level: tt.level,
			}
			actual := logConfig.GetSlogLevel()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
