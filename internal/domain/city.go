package domain

type City struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type GetAllCityCategoryResponse struct {
	Data     []*City        `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
