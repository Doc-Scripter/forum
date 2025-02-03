package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "forum/handlers"
	"forum/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler.Home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	fmt.Printf("Starting server on: %s", port)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)

	// Initialize the database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new user
	user := &database.User{Name: "John Doe"}
	if err := database.CreateUser(user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created user: %+v\n", user)

	// Retrieve the user
	retrievedUser, err := database.GetUser(user.ID)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	log.Printf("Retrieved user: %+v\n", retrievedUser)

	// Create a new product
	product := &database.Product{Name: "Laptop", Price: 999.99}
	if err := database.CreateProduct(product); err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	log.Printf("Created product: %+v\n", product)

	// Retrieve the product
	retrievedProduct, err := database.GetProduct(product.ID)
	if err != nil {
		log.Fatalf("Failed to get product: %v", err)
	}
	log.Printf("Retrieved product: %+v\n", retrievedProduct)
}
