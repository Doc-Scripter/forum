package models

import (
	e "forum/Error"
	d "forum/database"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ProfileData struct {
	Username string
	Email    string
	Uuid     string
	Initials string
	Category Categories
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
		Msg:  "The HTTP method used is not accessible",
		Code: http.StatusMethodNotAllowed,
	},
}

type Post struct {
	CreatedAt     time.Time `json:"created_at"`
	Categories    []string  `json:"category"`
	Likes         int       `json:"likes"`
	Title         string    `json:"title"`
	Dislikes      int       `json:"dislikes"`
	CommentsCount int       `json:"comments_count"`
	Comments      []Comment `json:"comments"`
	Content       string    `json:"content"`
	User_uuid     string    `json:"user_uuid"`
	Post_id       int       `json:"post_id"`
	Filepath      string    `json:"filepath"`
	Filename      string    `json:"filename"`
	Owner         string
	OwnerInitials string
}

type Categories struct {
	All_Categories string
	Technology string
	Health string
	Math string
	Games string
	Science string
	Religion string
	Education string
	Politics string
	Fashion string
	Lifestyle string
	Sports string
}

var Category = Categories{
	All_Categories : "All Categories",
	Technology :"Technology",
	Health :"Health",
	Math :"Math",
	Games :"Games",
	Science :"Science",
	Religion :"Religion",
	Education :"Education",
	Politics :"Politics",
	Fashion :"Fashion",
	Lifestyle :"Lifestyle",
	Sports : "Sports",
}
	

type Users struct {
	Username string
	Email    string
	Password string
}

var User Users

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_id    string    `json:"post_id"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Content    string    `json:"content"`
}

// Initials interface defines the method for generating initials
type Initials interface {
	GenerateInitials() string
}

// Make ProfileData implement the Initials interface
func (p *ProfileData) GenerateInitials() string {
	parts := strings.Split(p.Username, " ")

	if len(parts) == 0 || len(parts[0]) == 0 {
		return "?"
	}

	// Get first initial
	firstInitial := string(parts[0][0])

	// Get second initial if available
	var secondInitial string
	if len(parts) > 1 && len(parts[1]) > 0 {
		secondInitial = string(parts[1][0])
	}

	// Return combined initials
	if secondInitial != "" {
		return strings.ToUpper(firstInitial + secondInitial)
	}
	return strings.ToUpper(firstInitial)
}

type Image struct {
	ImageID   string    `json:"image_id"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}
type Category_Process interface {
	Seperate_Categories() string
}

// ===========The function will pack the categories as a slice of strings from the database==========
func (p *Post) Seperate_Categories() Post {
	var (
		combined_categories string
		categories          []string
	)

	row, err := d.Db.Query("SELECT category FROM posts")
	if err != nil {
		e.LogError(err)
		return *p
	}
	for row.Next() {
		err = row.Scan(&combined_categories)
		if err != nil {
			e.LogError(err)
			return Post{}
		}
		categories = strings.Split(combined_categories, ", ")
		p.Categories = categories
	}
	return *p
}

func (user *Users) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return nil
}


