package common

type Paging struct {
	Page  int   `JSON:"page" form:"page"`
	Limit int   `JSON:"limit" form:"limit"`
	Total int64 `JSON:"total" form:"total"`
}

func (p *Paging) Process() {
	if p.Page < 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}
