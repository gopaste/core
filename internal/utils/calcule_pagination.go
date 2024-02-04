package utils

import (
	"strconv"

	"github.com/Caixetadev/snippet/pkg/typesystem"
)

type Page struct {
	Page   int
	Offset int
	Limit  int
	Total  int
}

func CalculePagination(count int, pageStr string) (*Page, error) {
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, typesystem.ServerError
	}

	limit := 10

	totalPages := (count + limit - 1) / limit

	if totalPages < page {
		return nil, typesystem.NotFound
	}

	offset := (page - 1) * limit

	response := &Page{
		Total:  totalPages,
		Page:   page,
		Offset: offset,
		Limit:  limit,
	}

	return response, nil
}
