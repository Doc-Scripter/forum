package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Product represents a product in the database.
type Product struct {
	ID    string
	Name  string
	Price float64
}

// CreateProduct inserts a new product into the database.
func CreateProduct(product *Product) error {
	// Generate a new UUID for the product
	product.ID = uuid.New().String()

	query := `INSERT INTO products (id, name, price) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetProduct retrieves a product by ID.
func GetProduct(id string) (*Product, error) {
	product := &Product{}
	query := `SELECT id, name, price FROM products WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}