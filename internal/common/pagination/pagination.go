package pagination

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	PerPage     int `json:"per_page"`
}

type PaginationRequest struct {
	Page      int    `json:"page" query:"page" validate:"omitempty,gte=1"`
	Limit     int    `json:"limit" query:"limit" validate:"omitempty,gte=1,lte=100"`
	SortBy    string `json:"sortBy" query:"sortBy"`
	SortOrder string `json:"sortOrder" query:"sortOrder" validate:"omitempty,oneof=asc desc ASC DESC"`
	Search    string `json:"search" query:"search"`
}

func (r *PaginationRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

func BuildPagination(currentPage, totalData, perPage int) Pagination {
	totalPage := totalData / perPage
	if totalData%perPage != 0 {
		totalPage++
	}

	return Pagination{
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		TotalData:   totalData,
		PerPage:     perPage,
	}
}

func CalculateTotalPages(totalItems int64, itemsPerPage int) int {
	if itemsPerPage <= 0 {
		return 0
	}
	totalPages := totalItems / int64(itemsPerPage)
	if totalItems%int64(itemsPerPage) != 0 {
		totalPages++
	}
	return int(totalPages)
}
