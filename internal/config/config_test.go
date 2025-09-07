package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		description    string
	}{
		{
			name:    "default values when no env vars set",
			envVars: map[string]string{},
			expectedConfig: &Config{
				Port:      "8080",
				DBDialect: "sqlite",
				DBDSN:     "cupcake_store.db",
				LogLevel:  "info",
			},
			description: "should use default values when no environment variables are set",
		},
		{
			name: "environment variables override defaults",
			envVars: map[string]string{
				"PORT":       "9000",
				"DB_DIALECT": "postgres",
				"DB_DSN":     "host=test",
				"LOG_LEVEL":  "error",
			},
			expectedConfig: &Config{
				Port:      "9000",
				DBDialect: "postgres",
				DBDSN:     "host=test",
				LogLevel:  "error",
			},
			description: "should use environment variables when they are set",
		},
		{
			name: "partial environment variables",
			envVars: map[string]string{
				"PORT":   "9001",
				"DB_DSN": "host=partial",
			},
			expectedConfig: &Config{
				Port:      "9001",
				DBDialect: "sqlite",
				DBDSN:     "host=partial",
				LogLevel:  "info",
			},
			description: "should use defaults for missing environment variables",
		},
		{
			name: "empty environment variables use defaults",
			envVars: map[string]string{
				"PORT":       "",
				"DB_DIALECT": "",
				"DB_DSN":     "",
				"LOG_LEVEL":  "",
			},
			expectedConfig: &Config{
				Port:      "8080",
				DBDialect: "sqlite",
				DBDSN:     "cupcake_store.db",
				LogLevel:  "info",
			},
			description: "should use defaults when environment variables are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg := Load()

			require.Equal(t, tt.expectedConfig.Port, cfg.Port)
			require.Equal(t, tt.expectedConfig.DBDialect, cfg.DBDialect)
			require.Equal(t, tt.expectedConfig.DBDSN, cfg.DBDSN)
			require.Equal(t, tt.expectedConfig.LogLevel, cfg.LogLevel)

			os.Clearenv()
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
		description  string
	}{
		{
			name:         "returns environment value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "test_value",
			expected:     "test_value",
			description:  "should return environment variable value when it exists",
		},
		{
			name:         "returns default when env var not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
			description:  "should return default value when environment variable is not set",
		},
		{
			name:         "returns default when env var is empty",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
			description:  "should return default value when environment variable is empty",
		},
		{
			name:         "returns whitespace when env var is whitespace",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "   ",
			expected:     "   ",
			description:  "should return whitespace value when environment variable is only whitespace",
		},
		{
			name:         "returns newline when env var is newline",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "\n",
			expected:     "\n",
			description:  "should return newline value when environment variable is only newline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)

			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			val := getEnv(tt.key, tt.defaultValue)
			require.Equal(t, tt.expected, val)

			os.Unsetenv(tt.key)
		})
	}
}

func TestConfig_Fields(t *testing.T) {
	tests := []struct {
		name             string
		config           *Config
		expectedPort     string
		expectedDialect  string
		expectedDSN      string
		expectedLogLevel string
	}{
		{
			name: "all fields set",
			config: &Config{
				Port:      "8080",
				DBDialect: "sqlite",
				DBDSN:     "test.db",
				LogLevel:  "info",
			},
			expectedPort:     "8080",
			expectedDialect:  "sqlite",
			expectedDSN:      "test.db",
			expectedLogLevel: "info",
		},
		{
			name: "postgres configuration",
			config: &Config{
				Port:      "5432",
				DBDialect: "postgres",
				DBDSN:     "host=localhost user=postgres dbname=cupcake_store",
				LogLevel:  "debug",
			},
			expectedPort:     "5432",
			expectedDialect:  "postgres",
			expectedDSN:      "host=localhost user=postgres dbname=cupcake_store",
			expectedLogLevel: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedPort, tt.config.Port)
			require.Equal(t, tt.expectedDialect, tt.config.DBDialect)
			require.Equal(t, tt.expectedDSN, tt.config.DBDSN)
			require.Equal(t, tt.expectedLogLevel, tt.config.LogLevel)
		})
	}
}
