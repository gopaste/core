package pagination

import (
	"fmt"

	"github.com/Caixetadev/snippet/internal/entity"
	"github.com/Caixetadev/snippet/internal/utils"
)

func GeneratePaginationInfo(count int, page int, path string) (*entity.PaginationInfo, error) {
	limit := 10

	totalPages := (count + limit - 1) / limit

	nextPage, prevPage := "", ""
	if totalPages > page {
		nextPage = fmt.Sprintf("%s?page=%d", path, page+1)
	}

	if page > 1 {
		prevPage = fmt.Sprintf("%s?page=%d", path, page-1)
	}

	pagination := &entity.PaginationInfo{
		Pages: totalPages,
		Count: count,
		Next:  utils.StringToPtr(nextPage),
		Prev:  utils.StringToPtr(prevPage),
	}

	return pagination, nil
}
