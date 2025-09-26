// pkg/pagination/pagination.go
package pagination

type Query struct {
	Page  int    `query:"page"`  // 1-based
	Limit int    `query:"limit"` // per-page
	Sort  string `query:"sort"`  // e.g., "created_at:desc,name:asc"
	Q     string `query:"q"`     // free-text search
}

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

func (q *Query) Normalize() {
	if q.Page < 1 {
		q.Page = DefaultPage
	}
	if q.Limit <= 0 || q.Limit > MaxLimit {
		q.Limit = DefaultLimit
	}
}
