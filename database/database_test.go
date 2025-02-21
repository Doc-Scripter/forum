package database

import (
	"testing"
	"os"
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Modified to take a dbPath parameter for testing
func startDbConnectionWithPath(dbPath string) error {
	var err error
	Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	err = Db.Ping()
	if err != nil {
		return err
	}

	if err = CreateUsersTable(Db); err != nil {
		return err
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		return err
	}

	if err = CreateSessionsTable(Db); err != nil {
		return err
	}

	if err = CreatePostsTable(Db); err != nil {
		return err
	}
	if err = CreateCommentsTable(Db); err != nil {
		return err
	}
	
	return nil
}

func TestStartDbConnection(t *testing.T) {
	// Clean up any existing test database before running tests
	os.Remove("forum.db")

	tests := []struct {
		name    string
		wantErr bool
		dbPath  string
		setup   func(t *testing.T)
	}{
		{
			name:    "Successful connection and table creation",
			wantErr: false,
			dbPath:  "forum.db",
			setup:   func(t *testing.T) {},
		},
		{
			name:    "Invalid database path permission",
			wantErr: true,
			dbPath:  "", // Will be set in setup
			setup: func(t *testing.T) {
				// Create a read-only directory
				err := os.Mkdir("readonly", 0400)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
				// Use absolute path to ensure we hit the readonly directory
				absPath, err := filepath.Abs("./readonly/test.db")
				if err != nil {
					t.Fatalf("Setup failed getting absolute path: %v", err)
				}
				t.Setenv("TEST_DB_PATH", absPath)
			},
		},
		{
			name:    "Invalid database file",
			wantErr: true,
			dbPath:  "/dev/null/invalid.db", // Invalid path that should fail
			setup:   func(t *testing.T) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup
			if tt.setup != nil {
				tt.setup(t)
			}

			// Clean up database file before each test unless testing permissions
			if tt.name != "Invalid database path permission" {
				os.Remove("forum.db")
			}

			// Use the test-specific path if set, otherwise use default
			dbPath := tt.dbPath
			if tt.name == "Invalid database path permission" {
				dbPath = os.Getenv("TEST_DB_PATH")
			}

			// Reset Db to nil before each test
			if Db != nil {
				Db.Close()
				Db = nil
			}

			err := startDbConnectionWithPath(dbPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("startDbConnectionWithPath() error = %v, wantErr %v, path = %s", err, tt.wantErr, dbPath)
			}

			// Clean up after each test
			if tt.name == "Invalid database path permission" {
				os.RemoveAll("readonly")
			}
		})
	}

	// Final cleanup
	os.Remove("forum.db")
}