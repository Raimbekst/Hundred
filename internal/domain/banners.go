package domain

type Banner struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Status       int    `json:"status" db:"status"`
	Image        string `json:"image" db:"image"`
	Iframe       string `json:"iframe" db:"iframe"`
	LanguageType string `json:"language_type" db:"language_type"`
}

type GetAllBannersCategoryResponse struct {
	Data     []*Banner      `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
