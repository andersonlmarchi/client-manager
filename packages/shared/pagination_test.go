package shared_test

import (
	"testing"

	"github.com/andersonlmarchi/client-manager/packages/shared"
)

func TestPageQueryNormalizeAndOffset(t *testing.T) {
	t.Parallel()
	q := shared.PageQuery{Page: 0, PageSize: 0}.Normalize()
	if q.Page != 1 || q.PageSize != shared.DefaultPageSize {
		t.Fatalf("got %+v", q)
	}
	q = shared.PageQuery{Page: 2, PageSize: 1000}.Normalize()
	if q.PageSize != shared.MaxPageSize {
		t.Fatalf("page size = %d, want %d", q.PageSize, shared.MaxPageSize)
	}
	if off := (shared.PageQuery{Page: 3, PageSize: 10}).Offset(); off != 20 {
		t.Fatalf("offset = %d, want 20", off)
	}
	if lim := (shared.PageQuery{Page: 3, PageSize: 10}).Limit(); lim != 10 {
		t.Fatalf("limit = %d, want 10", lim)
	}
}

func TestNewPageResult(t *testing.T) {
	t.Parallel()
	res := shared.NewPageResult([]string{"a"}, 1, 10, 25)
	if res.TotalPages != 3 {
		t.Fatalf("TotalPages = %d, want 3", res.TotalPages)
	}
	res = shared.NewPageResult[string](nil, 1, 10, 0)
	if res.Items == nil || len(res.Items) != 0 {
		t.Fatalf("Items should be empty slice, got %#v", res.Items)
	}
}
