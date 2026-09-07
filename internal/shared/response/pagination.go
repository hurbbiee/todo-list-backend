package response

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Code       int64      `json:"code"`
	Message    string     `json:"message"`
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func Paginate[T any](
	data []T,
	page int64,
	limit int64,
	total int64,
) PaginatedResponse[T] {

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return PaginatedResponse[T]{
		Code:    0,
		Message: "success",
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
