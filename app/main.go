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

type Application struct {
	DB *sql.DB
}

func setUpApp() (Application, error) {
	var application Application
	// Connect to MySQL
	dsn := "root:password@tcp(db:3306)/mytheresa"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return application, fmt.Errorf("database connection failed: %v", err)
	}

	application.DB = db

	return application, nil
}

func (application *Application) productsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Products Handler starting...")
	defer log.Printf("Products Handler finished.")

	rows, err := application.DB.Query("SELECT sku, name, category, price FROM products LIMIT 5")
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
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func main() {
	application, err := setUpApp()
	if err != nil {
		log.Fatal("Failed to set up application:", err)
	}
	defer application.DB.Close()

	// Register handlers
	http.HandleFunc("/products", application.productsHandler)

	// Start server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
