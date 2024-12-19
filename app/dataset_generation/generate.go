package generator

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GenerateRows(db *sql.DB, totalProducts int, minPrice int, maxPrice int) {
	// Seed random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fetch last SKU
	var lastSKU int
	err := db.QueryRow("SELECT COALESCE(MAX(sku), 0) FROM products").Scan(&lastSKU)
	if err != nil {
		log.Fatalf("Failed to fetch last SKU: %v", err)
	}

	log.Printf("Starting product generation. Last SKU: %d", lastSKU)

	// Prepare for batch insertion
	batchSize := 1000
	products := make([]interface{}, 0, batchSize*4) // 4 fields per product
	query := "INSERT INTO products (sku, name, category, price) VALUES "

	// Generate and insert products
	for i := 1; i <= totalProducts; i++ {
		sku := fmt.Sprintf("%06d", lastSKU+i)
		name := fmt.Sprintf("Product %d", lastSKU+i)
		category := getCategory(lastSKU + i)
		price := rand.Intn(maxPrice-minPrice+1) + minPrice

		// Add product to batch
		query += "(?, ?, ?, ?),"
		products = append(products, sku, name, category, price)

		// Execute batch when it reaches the size
		if i%batchSize == 0 || i == totalProducts {
			query = query[:len(query)-1] // Remove trailing comma
			_, err := db.Exec(query, products...)
			if err != nil {
				log.Fatalf("Failed to execute batch insert: %v", err)
			}

			log.Printf("Inserted %d products", i)
			// Reset batch
			query = "INSERT INTO products (sku, name, category, price) VALUES "
			products = products[:0]
		}
	}

	log.Println("Finished inserting products")
}

// getCategory returns a category based on the product number
func getCategory(i int) string {
	switch i % 3 {
	case 0:
		return "boots"
	case 1:
		return "sandals"
	default:
		return "sneakers"
	}
}
