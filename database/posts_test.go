package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupPostsTestDB() (*sql.DB, error) {
	// Create a temporary database file
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func cleanupTestDB(db *sql.DB) {
	db.Close()
}

func TestCreatePostsTable(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Successful Table Creation",
			args:    args{db: func() *sql.DB { db, _ := setupPostsTestDB(); return db }()},
			wantErr: false,
		},
		{
			name:    "Nil Database Connection",
			args:    args{db: nil},
			wantErr: true,
		},
		{
			name:    "Table Already Exists",
			args:    args{db: func() *sql.DB { db, _ := setupTestDB(); return db }()},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Table Already Exists" {
				// First, create the table to simulate the scenario where the table already exists
				if err := CreatePostsTable(tt.args.db); err != nil {
					t.Fatalf("Failed to create table for 'Table Already Exists' test: %v", err)
				}
			}

			if err := CreatePostsTable(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreatePostsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		if tt.args.db != nil {
			cleanupTestDB(tt.args.db)
		}
	}
}
