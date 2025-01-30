package main

import (
	"time"

	m "forum/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Global variables
var (
	db    *gorm.DB
	store = session.New()
)

func main() {
	app := fiber.New()
	db = m.InitDB() // Initialize SQLite database

	// Authentication routes
	app.Get("/", Home)
	app.Post("/registration", RegisterUser)
	app.Post("/login", LoginUser)
	app.Get("/logout", LogoutUser)

	// Start the Fiber app
	app.Listen(":3000")
}

func Home(c *fiber.Ctx) error {
	return c.SendFile("templates/registration.html")
}

// Register a new user
func RegisterUser(c *fiber.Ctx) error {
	type Request struct {
		Email    string `form:"email"`
		Username string `form:"username"`
		Password string `form:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if email or username is already taken
	var existing m.User
	if err := db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email or username already taken"})
	}

	// Create new user
	user := m.User{
		UUID:     uuid.New().String(),
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	db.Create(&user)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login user and create session
func LoginUser(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Fetch user by email
	var user m.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create session
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	sess.Set("userID", user.UUID)
	sess.Set("username", user.Username)
	sess.Set("expires", time.Now().Add(24*time.Hour).Unix()) // Set session expiry

	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    user.UUID,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"message": "Login successful", "user": user.Username})
}

// Logout user
func LogoutUser(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session error"})
	}

	sess.Destroy()
	c.ClearCookie("session_id")

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// package main

// import (
// 	"html/template"
// 	"log"
// 	"net/http"
// )

// // Define home handler function which writes a byte slice
// func home(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" {
// 		http.Error(w, "page not found", 404)
// 	}
// 	t, err := template.ParseFiles("templates/home.html")
// 	if err != nil {
// 		log.Fatal("error parsing html")
// 	}
// 	t.Execute(w, r)
// }

// func main() {
// 	// Initialize a new servemux then register home function as handler for "/"
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", home)

// 	// Start a new web server
// 	// 8080 TCP network address to listen on
// 	log.Println("Starting server on: 8080")
// 	err := http.ListenAndServe(":8080", mux)
// 	log.Fatal(err)
// }
