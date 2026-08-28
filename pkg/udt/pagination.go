package udt

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" v:"required" comment:"页数"`
	PageSize int `json:"page_size" v:"required" comment:"页面数量"`
}

func NewPagination(page, pageSize int) *Pagination {
	return &Pagination{Page: page, PageSize: pageSize}
}

// From es
func (p Pagination) From() int {
	return (p.Page - 1) * p.PageSize
}

// Offset sqls
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Size ...
func (p Pagination) Size() int {
	return p.PageSize
}
