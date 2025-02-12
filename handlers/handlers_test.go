package handlers

import (
	"database/sql"
	"testing"

	"forum/database"
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

func TestGetCommentsByPostID(t *testing.T) {
	// Open an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create comments table: %v", err)
	}

	// Insert test data
	postID := 1
	comments := []struct {
		PostID  int
		Content string
	}{
		{PostID: postID, Content: "first comment"},
		{PostID: postID, Content: "second comment"},
		{PostID: 2, Content: "first comment"},
	}
	for _, comment := range comments {
		_, err := db.Exec("INSERT INTO comments (post_id, content) VALUES (?, ?)", comment.PostID, comment.Content)
		if err != nil {
			t.Fatalf("Failed to insert test comment: %v", err)
		}
	}

	// Retrieve comments for the post
	retrievedComments, err := GetCommentsByPostID(db, postID)
	if err != nil {
		t.Fatalf("GetCommentsByPostID failed: %v", err)
	}

	// Verify the number of comments retrieved
	expectedCount := 2 // Only 2 comments belong to postID 1
	if len(retrievedComments) != expectedCount {
		t.Errorf("Expected %d comments, got %d", expectedCount, len(retrievedComments))
	}

	// Verify the content of the retrieved comments
	for i, comment := range retrievedComments {
		if comment.PostID != postID {
			t.Errorf("Expected post ID %d, got %d", postID, comment.PostID)
		}
		if comment.Content != comments[i].Content {
			t.Errorf("Expected content %q, got %q", comments[i].Content, comment.Content)
		}
		if comment.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set, got zero value")
		}
	}
}
