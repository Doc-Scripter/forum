package models

import (
	"time"
)

type ProfileData struct {
	Username string
	Email    string
	Uuid     string
}

var Profile ProfileData

type Post struct {
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	Title     string    `json:"title"`
	Dislikes  int       `json:"dislikes"`
	Comments  string    `json:"comments"`
	Content   string    `json:"content"`
	User_uuid   string       `json:"user_uuid"`
	Post_id   int       `json:"post_id"`
}
