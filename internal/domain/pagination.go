package domain

type Pagination struct {
	Page  int `form:"page"  json:"page" validate:"min=0"`
	Limit int `form:"limit" json:"limit" validate:"min=0"`
}

type PaginationPage struct {
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Count int `json:"count"`
}
