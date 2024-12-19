package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryProductsFromDB(t *testing.T) {
	t.Run("Pagination limits to 5 products", func(t *testing.T) {
		// Create mock DB and set up rows
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000001", "Product 1", "boots", 100000).
			AddRow("000002", "Product 2", "boots", 120000).
			AddRow("000003", "Product 3", "sandals", 90000).
			AddRow("000004", "Product 4", "sneakers", 70000).
			AddRow("000005", "Product 5", "sneakers", 80000)

		// Update expected query to match the actual one
		mock.ExpectQuery(`SELECT sku, name, category, price FROM products LIMIT \? OFFSET \?`).
			WithArgs(5, 0).
			WillReturnRows(rows)

		// Call the function
		filters := Filters{}
		products, err := queryProductsFromDB(db, filters, 1, 10) // `10` to test if the function limits to `5`
		assert.NoError(t, err)

		// Verify results
		assert.Len(t, products, 5, "should return a maximum of 5 products")
		assert.Equal(t, "000005", products[4].SKU, "last product SKU should match")
	})

	t.Run("Max discount applied (SKU=000003 and category=boots)", func(t *testing.T) {
		// Create mock DB and set up rows
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Add a product with SKU=000003 and category=boots
		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000003", "Special Boots", "boots", 100000)

		// Update expected query to match the actual one
		mock.ExpectQuery(`SELECT sku, name, category, price FROM products LIMIT \? OFFSET \?`).
			WithArgs(1, 0). // Testing with a single row and no offset
			WillReturnRows(rows)

		// Call the function
		filters := Filters{}
		products, err := queryProductsFromDB(db, filters, 1, 1)
		assert.NoError(t, err)

		// Verify results
		assert.Len(t, products, 1, "should return one product")
		assert.Equal(t, "000003", products[0].SKU, "product SKU should match")
		assert.Equal(t, "boots", products[0].Category, "product category should match")
		assert.Equal(t, 70000, products[0].Price.Final, "final price should reflect 30% discount")
		assert.Equal(t, "30%", *products[0].Price.DiscountPercentage, "discount percentage should be 30%")
	})

	t.Run("Filter by category", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000001", "Boots 1", "boots", 100000).
			AddRow("000002", "Boots 2", "boots", 110000)

		mock.ExpectQuery("SELECT sku, name, category, price FROM products AND category = .* LIMIT .* OFFSET .*").
			WithArgs("boots", 5, 0).
			WillReturnRows(rows)

		filters := Filters{Category: "boots"}
		products, err := queryProductsFromDB(db, filters, 1, 5)
		assert.NoError(t, err)

		assert.Len(t, products, 2)
		assert.Equal(t, "boots", products[0].Category)
		assert.Equal(t, "boots", products[1].Category)
	})

	t.Run("Filter by price", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000004", "Affordable Sneakers", "sneakers", 70000)

		mock.ExpectQuery("SELECT sku, name, category, price FROM products AND price < .* LIMIT .* OFFSET .*").
			WithArgs(80000, 5, 0).
			WillReturnRows(rows)

		filters := Filters{PriceLessThan: 80000}
		products, err := queryProductsFromDB(db, filters, 1, 5)
		assert.NoError(t, err)

		assert.Len(t, products, 1)
		assert.Equal(t, "Affordable Sneakers", products[0].Name)
		assert.True(t, products[0].Price.Original < 80000)
	})

	t.Run("No filters applied", func(t *testing.T) {
		// Create mock DB and set up rows
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Add multiple rows to simulate products in the database
		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000001", "Product 1", "category1", 50000).
			AddRow("000002", "Product 2", "category2", 70000)

		// Update expected query to match the actual one
		mock.ExpectQuery(`SELECT sku, name, category, price FROM products LIMIT \? OFFSET \?`).
			WithArgs(5, 0). // Default pageSize=5, offset=0
			WillReturnRows(rows)

		// Call the function with no filters
		filters := Filters{}
		products, err := queryProductsFromDB(db, filters, 1, 5)
		assert.NoError(t, err)

		// Verify results
		assert.Len(t, products, 2, "should return all products")
		assert.Equal(t, "000001", products[0].SKU, "first product SKU should match")
		assert.Equal(t, "Product 1", products[0].Name, "first product name should match")
		assert.Equal(t, "000002", products[1].SKU, "second product SKU should match")
		assert.Equal(t, "Product 2", products[1].Name, "second product name should match")
	})

	t.Run("Different discounts applied on different products", func(t *testing.T) {
		// Create mock DB and set up rows
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Add two products to test different discount rules
		rows := sqlmock.NewRows([]string{"sku", "name", "category", "price"}).
			AddRow("000003", "Special Shoes", "shoes", 60000). // Should get 15% discount
			AddRow("000004", "Classic Boots", "boots", 80000)  // Should get 30% discount

		mock.ExpectQuery(`SELECT sku, name, category, price FROM products LIMIT \? OFFSET \?`).
			WithArgs(5, 0). // Default pageSize=5, offset=0
			WillReturnRows(rows)

		filters := Filters{}
		products, err := queryProductsFromDB(db, filters, 1, 5)
		assert.NoError(t, err)

		// Verify results
		assert.Len(t, products, 2, "should return two products")

		// Validate first product with SKU "000003"
		expectedFinalPrice1 := 51000 // 15% discount on 60000
		assert.Equal(t, "000003", products[0].SKU, "first product SKU should match")
		assert.Equal(t, "Special Shoes", products[0].Name, "first product name should match")
		assert.Equal(t, "shoes", products[0].Category, "first product category should match")
		assert.Equal(t, 60000, products[0].Price.Original, "first product original price should match")
		assert.Equal(t, expectedFinalPrice1, products[0].Price.Final, "first product final price should match")
		assert.Equal(t, "EUR", products[0].Price.Currency, "currency should be EUR")

		// Validate second product with category "boots"
		expectedFinalPrice2 := 56000 // 30% discount on 80000
		assert.Equal(t, "000004", products[1].SKU, "second product SKU should match")
		assert.Equal(t, "Classic Boots", products[1].Name, "second product name should match")
		assert.Equal(t, "boots", products[1].Category, "second product category should match")
		assert.Equal(t, 80000, products[1].Price.Original, "second product original price should match")
		assert.Equal(t, expectedFinalPrice2, products[1].Price.Final, "second product final price should match")
		assert.Equal(t, "EUR", products[1].Price.Currency, "currency should be EUR")
	})
}
