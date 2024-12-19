package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Application holds the application data, such as the database connection and HTTP server
type Application struct {
	DB     *sql.DB
	Server *http.Server
}

// SetUpApp initializes the database connection
func SetUpApp(dsn string) (*Application, error) {
	var application Application

	// Connect to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}

	application.DB = db

	mux := http.NewServeMux()

	// Register the products handler
	mux.HandleFunc("/products", application.ProductsHandler)

	// Set up the HTTP server
	application.Server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &application, nil
}

// ProductsHandler handles the fetching of products with optional pagination and filters
func (application *Application) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Products Handler starting...")
	defer log.Printf("Products Handler finished.")

	// Parse query parameters
	page, pageSize := parsePaginationParams(r)

	// Extract filters from query parameters (optional)
	filters := Filters{
		Category:      r.URL.Query().Get("category"),
		PriceLessThan: parseInt(r.URL.Query().Get("priceLessThan")),
	}

	// Fetch products with the specified filters and pagination
	products, err := GetProducts(application.DB, filters, page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting products: %v", err), http.StatusInternalServerError)
		return
	}

	// Write response
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
