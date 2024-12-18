package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
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
			if err := rows.Scan(&p.SKU, &p.Name, &p.Category, &p.Price); err != nil {
				http.Error(w, "Error scanning product", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}

		log.Println("Products fetched successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	// Start server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
