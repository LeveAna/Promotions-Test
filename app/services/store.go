package services

import (
	"database/sql"
	"fmt"
	"strconv"
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

// TODO: move to database for more flexibility
type Discount struct {
	Percentage int
	SKU        *string
	Category   *string
}

type Filters struct {
	Category      string
	PriceLessThan int
}

// TODO: add caching

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

		// Apply discounts
		discounts := []Discount{
			// Rule 1: 30% discount for "boots" category
			{Percentage: 30, Category: strPtr("boots")},
			// Rule 2: 15% discount for SKU "000003"
			{Percentage: 15, SKU: strPtr("000003")},
		}

		applyDiscount(&p, discounts)

		p.Price.Currency = "EUR"
		products = append(products, p)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %v", err)
	}

	return products, nil
}

// applyDiscount applies the discount rules to a product
func applyDiscount(p *Product, discounts []Discount) {
	if p == nil {
		return
	}

	// Default final price is the original price
	p.Price.Final = p.Price.Original

	var maxDiscount Discount
	for _, d := range discounts {
		if (d.Category != nil && *d.Category == p.Category) || (d.SKU != nil && *d.SKU == p.SKU) {
			if d.Percentage > maxDiscount.Percentage {
				maxDiscount = d
			}
		}
	}

	if maxDiscount.Percentage > 0 {
		discountAmount := p.Price.Original * maxDiscount.Percentage / 100
		p.Price.Final = p.Price.Original - discountAmount

		discountLabel := strconv.Itoa(maxDiscount.Percentage) + "%"
		p.Price.DiscountPercentage = &discountLabel
	}
}
