package models

import (
	"net/http"
	"time"
)

type ProfileData struct {
	Username string
	Email    string
	Uuid     string
}

type ErrorData struct {
	Msg  string
	Code int
}

type Errors struct {
	InternalError    ErrorData
	PageNotFound     ErrorData
	BadRequest       ErrorData
	Unauthorized     ErrorData
	Forbidden        ErrorData
	MethodNotAllowed ErrorData
}

var ErrorsData = Errors{
	InternalError: ErrorData{
		Msg:  "An unexpected error occurred. The error seems to be on our end. Hang tight",
		Code: http.StatusInternalServerError,
	},
	PageNotFound: ErrorData{
		Msg:  "The Page you are trying to access does not seem to exist",
		Code: http.StatusNotFound,
	},
	BadRequest: ErrorData{
		Msg:  "That's a Bad Request",
		Code: http.StatusBadRequest,
	},

	Unauthorized: ErrorData{
		Msg:  "You are not authorized to view this page",
		Code: http.StatusUnauthorized,
	},

	Forbidden: ErrorData{
		Msg:  "You are not allowed to perform this action",
		Code: http.StatusForbidden,
	},

	MethodNotAllowed: ErrorData{
		Msg:  "The HTTP method you used is not allowed for this route",
		Code: http.StatusMethodNotAllowed,
	},
}


type Post struct {
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	Title     string    `json:"title"`
	Dislikes  int       `json:"dislikes"`
	CommentsCount  int   `json:"comments_count"`
	Comments  []string `json:"comments"`
	Content   string    `json:"content"`
	User_uuid string    `json:"user_uuid"`
	Post_id   int       `json:"post_id"`
}

type Comment struct{
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Content   string    `json:"content"`

}
type Image struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    PostID    int64     `json:"post_id,omitempty"`
    Filename  string    `json:"filename"`
    Path      string    `json:"path"`
    CreatedAt time.Time `json:"created_at"`
}
