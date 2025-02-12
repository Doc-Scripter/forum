package handlers

import (
	"database/sql"
	"forum/database"
	"testing"
)

// TestAddComment tests the AddComment function
func TestAddComment(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	// Create the comments table
	if err := database.CreateCommentsTable(db); err != nil {
		t.Fatalf("Failed to create comments table: %v", err)
	}
	// Test data
	postID := 1
	content := "Test comment"

	// Add a comment
	commentID, err := AddComment(db, postID, content)
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}
	if commentID <= 0 {
		t.Errorf("Expected comment ID to be greater than 0, got %d", commentID)
	}
	// Query the database to verify the comment
	var retrievedContent string
	var retrievedPostID int
	err = db.QueryRow("SELECT post_id, content FROM comments WHERE comment_id = ?", commentID).Scan(&retrievedPostID, &retrievedContent)
	if err != nil {
		t.Fatalf("failed to query comment: %v", err)
	}

	// Verify the retrieved data
	if retrievedPostID != postID {
		t.Errorf("Expected post ID %d, got %d", postID, retrievedPostID)
	}
	if retrievedContent != content {
		t.Errorf("Expected content %q, got %q", content, retrievedContent)
	}
}
