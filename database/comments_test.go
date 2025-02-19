package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, error) {
	// Create a temporary database file for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateCommentsTable(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful Table Creation",
			args: args{
				db: func() *sql.DB {
					db, err := setupTestDB()
					if err != nil {
						t.Fatalf("Failed to setup test database: %v", err)
					}
					return db
				}(),
			},
			wantErr: false,
		},
		{
			name: "Database Connection Failure",
			args: args{
				db: nil,
			},
			wantErr: true,
		},
		{
			name: "Table Already Exists",
			args: args{
				db: func() *sql.DB {
					db, err := setupTestDB()
					if err != nil {
						t.Fatalf("Failed to setup test database: %v", err)
					}
					// Create the table first
					if err := CreateCommentsTable(db); err != nil {
						t.Fatalf("Failed to create table: %v", err)
					}
					return db
				}(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateCommentsTable(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateCommentsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}