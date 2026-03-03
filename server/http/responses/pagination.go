package responses

import "math"

type PaginationMeta struct {
	PerPage      uint `json:"per_page"`
	Total        uint `json:"total"`
	CurrentPage  uint `json:"current_page"`
	LastPage     uint `json:"last_page"`
	HasMorePages bool `json:"has_more_pages"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

const DefaultPerPage = 10
const MaxPerPage = 100

func NewPaginatedResponse(data interface{}, meta PaginationMeta) PaginatedResponse {
	return PaginatedResponse{Data: data, Meta: meta}
}

func NewPaginationMeta(page, perPage, total uint) PaginationMeta {
	lastPage := uint(math.Max(float64(total/perPage), 1))

	return PaginationMeta{
		PerPage:      perPage,
		Total:        total,
		CurrentPage:  page,
		LastPage:     lastPage,
		HasMorePages: page > lastPage,
	}
}
