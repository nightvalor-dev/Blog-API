package utils

import (
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

func ParsePagination(r *http.Request) PaginationParams {
	page := ParseIntParam(r, "page", 1)
	limit := ParseIntParam(r, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func ParseIntParam(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func ParseIntQuery(r *http.Request, key string) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return 0
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return n
}
