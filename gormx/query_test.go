package gormx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type pageData struct {
	Page     uint64
	PageSize uint64
	OrderBy  string
}

func (p *pageData) GetOrderBy() string {
	return p.OrderBy
}

func (p *pageData) GetPage() uint64 {
	return p.Page
}

func (p *pageData) GetPageSize() uint64 {
	return p.PageSize
}

func TestQuery_Page(t *testing.T) {
	tests := []struct {
		name string
		page Page
		want *Query
	}{
		{
			name: "nil",
			page: nil,
			want: NewQuery(),
		},
		{
			name: "only orderby",
			page: &pageData{
				OrderBy: "{\"id\":\"desc\"}",
			},
			want: NewQuery().OrderBy("id desc"),
		},
		{
			name: "only pagenum and size",
			page: &pageData{
				Page:     2,
				PageSize: 10,
			},
			want: NewQuery().Limit(10).Offset(10),
		},
		{
			name: "all params",
			page: &pageData{
				Page:     2,
				PageSize: 10,
				OrderBy:  "{\"id\":\"desc\"}",
			},
			want: NewQuery().Limit(10).Offset(10).OrderBy("id desc"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQuery().Page(tt.page)
			assert.Equal(t, tt.want, got)
		})
	}
}
