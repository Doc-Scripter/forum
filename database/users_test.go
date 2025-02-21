package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateUsersTable(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid SQLite database",
			args: args{
				// Create an in-memory SQLite database for testing
				db: func() *sql.DB {
					db, err := sql.Open("sqlite3", ":memory:")
					if err != nil {
						t.Fatalf("Failed to create test database: %v", err)
					}
					return db
				}(),
			},
			wantErr: false,
		},
		{
			name: "Nil database connection",
			args: args{
				db: nil,
			},
			wantErr: true,
		},
		{
			name: "Closed database connection",
			args: args{
				db: func() *sql.DB {
					db, err := sql.Open("sqlite3", ":memory:")
					if err != nil {
						t.Fatalf("Failed to create test database: %v", err)
					}
					db.Close()
					return db
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateUsersTable(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateUsersTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Cleanup: close database if it was opened
			if tt.args.db != nil && !tt.wantErr {
				tt.args.db.Close()
			}
		})
	}
}