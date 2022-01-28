package domain

type Faq struct {
	Id           int    `json:"id,omitempty" db:"id"`
	Question     string `json:"question" db:"question"`
	Answer       string `json:"answer" db:"answer"`
	LanguageType string `json:"language_type" db:"language_type"`
}
type GetAllFaqsCategoryResponse struct {
	Data     []*Faq         `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}

type Description struct {
	Id           int    `json:"id,omitempty" db:"id"`
	Caption      string `json:"caption" db:"caption"`
	Text         string `json:"text" db:"text"`
	LanguageType string `json:"language_type" db:"language_type"`
}

type GetAllDescCategoryResponse struct {
	Data     []*Description `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
