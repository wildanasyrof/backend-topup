// pkg/pagination/response.go
package pagination

type Page[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

type Meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`      // total rows (with filters)
	TotalPage int `json:"total_page"` // ceil(total/limit)
}
