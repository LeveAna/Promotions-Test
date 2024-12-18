package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	// Connect to MySQL
	dsn := "root:password@tcp(db:3306)/mytheresa"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Handle products endpoint
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT sku, name, category, price FROM products LIMIT 5")
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.SKU, &p.Name, &p.Category, &p.Price.Original); err != nil {
				http.Error(w, "Error scanning product", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		log.Println("Products fetched successfully")
		w.Header().Set("Content-Type", "application/json")

		// Marshal the products slice into pretty JSON format
		formattedJSON, err := json.MarshalIndent(products, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write(formattedJSON)
	})

	// Start server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
