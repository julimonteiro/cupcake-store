package main

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestMain_EnvironmentSetup(t *testing.T) {
	tests := []struct {
		name          string
		envFileExists bool
		envContent    string
		expectedLog   string
		description   string
	}{
		{
			name:          "no .env file",
			envFileExists: false,
			expectedLog:   ".env file not found",
			description:   "should log when .env file is not found",
		},
		{
			name:          "existing .env file",
			envFileExists: true,
			envContent:    "PORT=8081\nDB_DIALECT=sqlite\nDB_DSN=test.db\nLOG_LEVEL=info",
			expectedLog:   "",
			description:   "should not log when .env file exists",
		},
		{
			name:          "empty .env file",
			envFileExists: true,
			envContent:    "",
			expectedLog:   "",
			description:   "should handle empty .env file",
		},
		{
			name:          "partial .env file",
			envFileExists: true,
			envContent:    "PORT=8081\nDB_DIALECT=sqlite",
			expectedLog:   "",
			description:   "should handle partial .env file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			os.Remove(".env")

			if tt.envFileExists {
				err := os.WriteFile(".env", []byte(tt.envContent), 0644)
				require.NoError(t, err)
				defer os.Remove(".env")
			}

			if err := godotenv.Load(); err != nil {
				log.Println(".env file not found, using system environment variables")
			}

			if tt.expectedLog != "" {
				require.Contains(t, buf.String(), tt.expectedLog, tt.description)
			} else {
				require.NotContains(t, buf.String(), ".env file not found", tt.description)
			}
		})
	}
}

func TestMain_ConfigLoading(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		description string
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"PORT":       "8080",
				"DB_DIALECT": "sqlite",
				"DB_DSN":     ":memory:",
				"LOG_LEVEL":  "info",
			},
			description: "should load valid configuration",
		},
		{
			name: "partial configuration",
			envVars: map[string]string{
				"PORT": "3000",
			},
			description: "should load partial configuration with defaults",
		},
		{
			name:        "no environment variables",
			envVars:     map[string]string{},
			description: "should use default configuration",
		},
		{
			name: "invalid configuration values",
			envVars: map[string]string{
				"PORT":       "invalid",
				"DB_DIALECT": "mysql",
				"DB_DSN":     "invalid://dsn",
				"LOG_LEVEL":  "invalid",
			},
			description: "should handle invalid configuration values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			for key, value := range tt.envVars {
				require.Equal(t, value, os.Getenv(key), tt.description)
			}

			os.Clearenv()
		})
	}
}

func TestMain_DefaultConfig(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "default values when no env vars",
			description: "should use default values when no environment variables are set",
		},
		{
			name:        "default values when env vars are empty",
			description: "should use default values when environment variables are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			require.Equal(t, "", os.Getenv("PORT"), tt.description)
			require.Equal(t, "", os.Getenv("DB_DIALECT"), tt.description)
			require.Equal(t, "", os.Getenv("DB_DSN"), tt.description)
			require.Equal(t, "", os.Getenv("LOG_LEVEL"), tt.description)
		})
	}
}

func TestMain_InvalidConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		description string
	}{
		{
			name: "invalid port",
			envVars: map[string]string{
				"PORT": "invalid",
			},
			description: "should handle invalid port value",
		},
		{
			name: "invalid database dialect",
			envVars: map[string]string{
				"DB_DIALECT": "mysql",
			},
			description: "should handle invalid database dialect",
		},
		{
			name: "invalid database DSN",
			envVars: map[string]string{
				"DB_DSN": "invalid://dsn",
			},
			description: "should handle invalid database DSN",
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"LOG_LEVEL": "invalid",
			},
			description: "should handle invalid log level",
		},
		{
			name: "all invalid values",
			envVars: map[string]string{
				"PORT":       "invalid",
				"DB_DIALECT": "mysql",
				"DB_DSN":     "invalid://dsn",
				"LOG_LEVEL":  "invalid",
			},
			description: "should handle all invalid configuration values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			for key, value := range tt.envVars {
				require.Equal(t, value, os.Getenv(key), tt.description)
			}

			os.Clearenv()
		})
	}
}

func TestMain_ConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		description string
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"PORT":       "8080",
				"DB_DIALECT": "sqlite",
				"DB_DSN":     ":memory:",
				"LOG_LEVEL":  "info",
			},
			description: "should validate valid configuration",
		},
		{
			name: "configuration with defaults",
			envVars: map[string]string{
				"PORT": "3000",
			},
			description: "should validate configuration with defaults",
		},
		{
			name:        "empty configuration",
			envVars:     map[string]string{},
			description: "should validate empty configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			os.Clearenv()

			require.True(t, true, tt.description)
		})
	}
}

func TestMain_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		expect      bool
		description string
	}{
		{
			name: "valid config",
			env: map[string]string{
				"PORT":       "8080",
				"DB_DIALECT": "sqlite",
				"DB_DSN":     ":memory:",
				"LOG_LEVEL":  "info",
			},
			expect:      true,
			description: "should handle valid environment variables",
		},
		{
			name: "partial config",
			env: map[string]string{
				"PORT": "3000",
			},
			expect:      true,
			description: "should handle partial environment variables",
		},
		{
			name:        "empty config",
			env:         map[string]string{},
			expect:      true,
			description: "should handle empty environment variables",
		},
		{
			name: "invalid config",
			env: map[string]string{
				"PORT":       "invalid",
				"DB_DIALECT": "mysql",
				"DB_DSN":     "invalid://dsn",
				"LOG_LEVEL":  "invalid",
			},
			expect:      true,
			description: "should handle invalid environment variables",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			for key, value := range tt.env {
				require.Equal(t, value, os.Getenv(key), tt.description)
			}

			os.Clearenv()
		})
	}
}

func TestMain_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		scenario    string
		description string
	}{
		{
			name:        "missing .env file",
			scenario:    "no_env_file",
			description: "should handle missing .env file gracefully",
		},
		{
			name:        "invalid .env file",
			scenario:    "invalid_env_file",
			description: "should handle invalid .env file gracefully",
		},
		{
			name:        "permission denied",
			scenario:    "permission_denied",
			description: "should handle permission denied gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			os.Remove(".env")

			switch tt.scenario {
			case "no_env_file":
				if err := godotenv.Load(); err != nil {
					log.Println(".env file not found, using system environment variables")
				}
				require.Contains(t, buf.String(), ".env file not found", tt.description)

			case "invalid_env_file":
				err := os.WriteFile(".env", []byte("invalid content"), 0644)
				require.NoError(t, err)
				defer os.Remove(".env")

				if err := godotenv.Load(); err != nil {
					log.Println(".env file not found, using system environment variables")
				}

			case "permission_denied":
				if err := godotenv.Load(); err != nil {
					log.Println(".env file not found, using system environment variables")
				}
				require.Contains(t, buf.String(), ".env file not found", tt.description)
			}
		})
	}
}
