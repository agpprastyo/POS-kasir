package pagination

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalData   int `json:"total_data"`
	PerPage     int `json:"per_page"`
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
