package shared

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// PageQuery holds offset pagination inputs.
type PageQuery struct {
	Page     int
	PageSize int
}

// PageResult is a page of items with totals.
type PageResult[T any] struct {
	Items      []T
	Page       int
	PageSize   int
	TotalItems int64
	TotalPages int
}

// Normalize applies defaults and clamps page/page_size to safe bounds.
func (q PageQuery) Normalize() PageQuery {
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = DefaultPageSize
	}
	if size > MaxPageSize {
		size = MaxPageSize
	}
	return PageQuery{Page: page, PageSize: size}
}

// Offset returns the SQL OFFSET for the normalized query.
func (q PageQuery) Offset() int {
	n := q.Normalize()
	return (n.Page - 1) * n.PageSize
}

// Limit returns the SQL LIMIT for the normalized query.
func (q PageQuery) Limit() int {
	return q.Normalize().PageSize
}

// NewPageResult builds a PageResult and computes TotalPages.
func NewPageResult[T any](items []T, page, pageSize int, totalItems int64) PageResult[T] {
	if items == nil {
		items = []T{}
	}
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	}
	return PageResult[T]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
