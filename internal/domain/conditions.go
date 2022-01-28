package domain

type Condition struct {
	Id           int    `json:"id,omitempty" db:"id"`
	Caption      string `json:"caption" db:"caption"`
	Text         string `json:"text" db:"text"`
	LanguageType string `json:"language_type" db:"language_type"`
}

type GetAllConditionCategoryResponse struct {
	Data     []*Condition   `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
