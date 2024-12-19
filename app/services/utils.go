package services

import (
	"net/http"
	"strconv"
)

// strPtr returns a pointer to the given string
func strPtr(s string) *string {
	return &s
}

// parsePaginationParams extracts pagination parameters from the query string
func parsePaginationParams(r *http.Request) (int, int) {
	// Default values
	defaultPage := 1
	defaultPageSize := 10

	// Parse "page"
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	// Parse "pageSize"
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	return page, pageSize
}

// parseInt safely parses a query parameter to an integer, returning 0 if invalid
func parseInt(param string) int {
	if param == "" {
		return 0
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return value
}
