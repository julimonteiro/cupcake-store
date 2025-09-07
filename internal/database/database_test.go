package database

import (
	"testing"

	"github.com/julimonteiro/cupcake-store/internal/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectedError  string
		validateResult func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "SQLite with in-memory database",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "error",
			},
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name: "SQLite with file database",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     "test.db",
				LogLevel:  "error",
			},
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name: "PostgreSQL connection (expected to fail)",
			config: &config.Config{
				DBDialect: "postgres",
				DBDSN:     "postgres://user:pass@localhost:5432/test?sslmode=disable",
				LogLevel:  "error",
			},
			expectedError: "error connecting to database",
		},
		{
			name: "unsupported database dialect",
			config: &config.Config{
				DBDialect: "unsupported",
				DBDSN:     "test.db",
				LogLevel:  "info",
			},
			expectedError: "unsupported database dialect",
		},
		{
			name: "invalid DSN",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     "invalid://dsn",
				LogLevel:  "error",
			},
			expectedError: "error connecting to database",
		},
		{
			name: "MySQL connection (expected to fail)",
			config: &config.Config{
				DBDialect: "mysql",
				DBDSN:     "user:pass@tcp(localhost:3306)/test",
				LogLevel:  "error",
			},
			expectedError: "unsupported database dialect",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Init(tt.config)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, db)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				if tt.validateResult != nil {
					tt.validateResult(t, db)
				}
			}
		})
	}
}

func TestInit_LogLevels(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       string
		expectedError  string
		validateResult func(t *testing.T, db *gorm.DB)
	}{
		{
			name:     "info log level",
			logLevel: "info",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:     "error log level",
			logLevel: "error",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:     "debug log level",
			logLevel: "debug",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:     "warn log level",
			logLevel: "warn",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:     "silent log level",
			logLevel: "silent",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:     "invalid log level (defaults to error)",
			logLevel: "invalid",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  tt.logLevel,
			}

			db, err := Init(cfg)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, db)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				if tt.validateResult != nil {
					tt.validateResult(t, db)
				}
			}
		})
	}
}

func TestRunMigrations(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectedError  string
		validateResult func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "migrations run successfully",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "error",
			},
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)

				var count int64
				err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='cupcakes'").Scan(&count).Error
				require.NoError(t, err)
				require.Equal(t, int64(1), count)

				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name: "migrations with info log level",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "info",
			},
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)

				var count int64
				err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='cupcakes'").Scan(&count).Error
				require.NoError(t, err)
				require.Equal(t, int64(1), count)

				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name: "migrations with debug log level",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     ":memory:",
				LogLevel:  "debug",
			},
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)

				var count int64
				err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='cupcakes'").Scan(&count).Error
				require.NoError(t, err)
				require.Equal(t, int64(1), count)

				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Init(tt.config)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, db)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				if tt.validateResult != nil {
					tt.validateResult(t, db)
				}
			}
		})
	}
}

func TestInit_DatabaseTypes(t *testing.T) {
	tests := []struct {
		name           string
		dialect        string
		dsn            string
		expectedError  string
		validateResult func(t *testing.T, db *gorm.DB)
	}{
		{
			name:    "SQLite with memory database",
			dialect: "sqlite",
			dsn:     ":memory:",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:    "SQLite with file database",
			dialect: "sqlite",
			dsn:     "test.db",
			validateResult: func(t *testing.T, db *gorm.DB) {
				require.NotNil(t, db)
				sqlDB, err := db.DB()
				require.NoError(t, err)
				require.NoError(t, sqlDB.Close())
			},
		},
		{
			name:          "PostgreSQL (expected to fail)",
			dialect:       "postgres",
			dsn:           "postgres://user:pass@localhost:5432/test?sslmode=disable",
			expectedError: "error connecting to database",
		},
		{
			name:          "MySQL (unsupported)",
			dialect:       "mysql",
			dsn:           "user:pass@tcp(localhost:3306)/test",
			expectedError: "unsupported database dialect",
		},
		{
			name:          "SQL Server (unsupported)",
			dialect:       "mssql",
			dsn:           "sqlserver://user:pass@localhost:1433?database=test",
			expectedError: "unsupported database dialect",
		},
		{
			name:          "Oracle (unsupported)",
			dialect:       "oracle",
			dsn:           "oracle://user:pass@localhost:1521/test",
			expectedError: "unsupported database dialect",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DBDialect: tt.dialect,
				DBDSN:     tt.dsn,
				LogLevel:  "error",
			}

			db, err := Init(cfg)

			if tt.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, db)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, db)
				if tt.validateResult != nil {
					tt.validateResult(t, db)
				}
			}
		})
	}
}

func TestInit_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.Config
		expectedError string
	}{
		{
			name: "invalid SQLite DSN",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     "invalid://dsn",
				LogLevel:  "error",
			},
			expectedError: "error connecting to database",
		},
		{
			name: "invalid PostgreSQL DSN",
			config: &config.Config{
				DBDialect: "postgres",
				DBDSN:     "invalid://dsn",
				LogLevel:  "error",
			},
			expectedError: "error connecting to database",
		},
		{
			name: "empty dialect",
			config: &config.Config{
				DBDialect: "",
				DBDSN:     ":memory:",
				LogLevel:  "error",
			},
			expectedError: "unsupported database dialect",
		},
		{
			name: "empty DSN",
			config: &config.Config{
				DBDialect: "sqlite",
				DBDSN:     "",
				LogLevel:  "error",
			},
			expectedError: "error connecting to database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Init(tt.config)

			require.Error(t, err)
			require.Nil(t, db)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
