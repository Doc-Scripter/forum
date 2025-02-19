package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func cleanupDB(db *sql.DB) {
	db.Close()
}

type args struct {
	db *sql.DB
}

func TestCreateLikesDislikesTable(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer cleanupDB(db)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful Table Creation",
			args: args{
				db: db,
			},
			wantErr: false,
		},
		{
			name: "Table Already Exists",
			args: args{
				db: db,
			},
			wantErr: false,
		},
		{
			name: "Nil Database Connection",
			args: args{
				db: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Table Already Exists" {
				// Create the table first to simulate the scenario where the table already exists
				err := CreateLikesDislikesTable(tt.args.db)
				if err != nil {
					t.Errorf("Failed to create table for 'Table Already Exists' test case: %v", err)
				}
			}

			err := CreateLikesDislikesTable(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLikesDislikesTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
