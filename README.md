# Promotions-Test

A technical challenge for job selection, designed to demonstrate skills in backend development, database integration, caching, and API design.

## Features:
- List products with pagination and filters (e.g., category, price).
- Redis-based caching for better performance (optional).
- Apply predefined discounts based on product category or SKU.

## Prerequisites:
To run this project, you'll need the following installed:

- **Go (1.22 or higher)**
- **MySQL 8.0**
- **Redis** (if caching is enabled)

## Setup:

The application is set up with Docker containers. You can start the entire application (including MySQL and Redis) using Docker Compose.

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/promotions-test.git
    cd promotions-test
    ```

2. Start the services:
    ```bash
    docker-compose up
    ```

This will set up the MySQL database, Redis container, and the application itself. \
After running the command, please allow 5 seconds for the database to be ready before the application starts. When database is ready, `dataset_generator` generates up to 30000 database entries and saves it to database.

## Running Tests:

To run the tests, use the following command:

```bash
go test -v ./...
```

### API endpoint

#### GET `/products`
Returns a list of products.

**Query Parameters:**
  * `page` (optional): Page number (default is 1)
  * `pageSize` (optional): Number of products per page (default and maximum is 5)
  * `category` (optional): Filter by product category
  * `priceLessThan` (optional): Filter by price less than a specified amount.

### Example Request and Response

#### Request

**GET** `/products?page=1&pageSize=5&category=boots&priceLessThan=100000`

This request fetches a list of products with the following filters:
- Page 1
- 5 products per page (default and maximum)
- Category filter set to `boots`
- Price less than 1000.00 EUR

#### Response
```json
[
  {
    "sku": "000001",
    "name": "BV Lean leather ankle boots",
    "category": "boots",
    "price": {
      "original": 89000,
      "final": 89000,
      "currency": "EUR"
    }
  },
  {
    "sku": "000002",
    "name": "Nike Air Zoom Pegasus",
    "category": "boots",
    "price": {
      "original": 95000,
      "final": 95000,
      "currency": "EUR"
    }
  }
]
```

### Code structure
* Source code is located in `app` folder. `main.go` contains `main` function that starts the application.
* `dataset_generation/` contains a method for generating additional database entries and saving them to database.
* `services/` contains handler method, store method for fetching products and utility methods.
