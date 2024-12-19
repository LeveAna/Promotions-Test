package services

import (
	"database/sql"
	"fmt"
)

type Product struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    Price  `json:"price"`
}

type Price struct {
	Original           int     `json:"original"`
	Final              int     `json:"final"`
	DiscountPercentage *string `json:"discount_percentage,omitempty"`
	Currency           string  `json:"currency"`
}

type Filters struct {
	Category      string
	PriceLessThan int
}

// GetProducts retrieves products from the database with optional filters and simple pagination.
func GetProducts(db *sql.DB, filters Filters, page, pageSize int) ([]Product, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 5 {
		pageSize = 5 // Return up to 5 products per page
	}

	// Start building the SQL query
	query := "SELECT sku, name, category, price FROM products WHERE 1=1"
	args := []interface{}{} // To hold query parameters

	// Apply category filter if it's provided
	if filters.Category != "" {
		query += " AND category = ?"
		args = append(args, filters.Category)
	}

	// Apply price filter if it's provided
	if filters.PriceLessThan > 0 {
		query += " AND price < ?"
		args = append(args, filters.PriceLessThan)
	}

	// Add pagination
	offset := (page - 1) * pageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	// Execute the query with the dynamic conditions
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.SKU, &p.Name, &p.Category, &p.Price.Original); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, p)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %v", err)
	}

	return products, nil
}
