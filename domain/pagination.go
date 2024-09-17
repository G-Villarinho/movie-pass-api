package domain

import "strconv"

type Pagination struct {
	Limit      int    `json:"limit,omitempty" query:"limit"`
	Page       int    `json:"page,omitempty" query:"page"`
	Sort       string `json:"sort,omitempty" query:"sort"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
	Rows       any    `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) SetLimit(limit string) {
	if l, err := strconv.Atoi(limit); err == nil {
		p.Limit = l
	} else {
		p.Limit = 10
	}
}

func (p *Pagination) SetPage(page string) {
	if pg, err := strconv.Atoi(page); err == nil {
		p.Page = pg
	} else {
		p.Page = 1
	}
}

func (p *Pagination) SetSort(sort string) {
	if sort != "" {
		p.Sort = sort
	}
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}
