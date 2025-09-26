// pkg/pagination/scopes.go
package pagination

import (
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ScopePaginate(db *gorm.DB, page, limit int) *gorm.DB {
	offset := (page - 1) * limit
	return db.Offset(offset).Limit(limit)
}

func ScopeSort(db *gorm.DB, sort string, allowed map[string]struct{}) *gorm.DB {
	if sort == "" { // default
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
	}
	parts := strings.Split(sort, ",")
	for _, p := range parts {
		seg := strings.SplitN(strings.TrimSpace(p), ":", 2)
		col := seg[0]
		dir := "asc"
		if len(seg) == 2 && strings.EqualFold(seg[1], "desc") {
			dir = "desc"
		}
		if _, ok := allowed[col]; ok {
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: col}, Desc: dir == "desc"})
		}
	}
	return db
}

func CalcMeta(total, page, limit int) Meta {
	tp := int(math.Ceil(float64(total) / float64(limit)))
	return Meta{Page: page, Limit: limit, Total: total, TotalPage: tp}
}
