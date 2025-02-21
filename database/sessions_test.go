package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateSessionsTable(t *testing.T) {
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
				db: func() *sql.DB {
					db, _ := sql.Open("sqlite3", ":memory:") // In-memory database for testing
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateSessionsTable(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("CreateSessionsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Clean up: close the database if it was opened
			if tt.args.db != nil {
				tt.args.db.Close()
			}
		})
	}
}
