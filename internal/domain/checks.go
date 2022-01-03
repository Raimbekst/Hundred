package domain

type CheckInfo struct {
	Id           int          `json:"id" db:"id"`
	UserId       int          `json:"user_id" db:"user_id"`
	PartnerId    int          `json:"partner_id" db:"partner_id"`
	CheckAmount  int          `json:"check_amount" db:"check_amount"`
	CheckDate    float64      `json:"check_date" db:"check_date"`
	IsWinner     bool         `json:"is_winner" db:"is_winner"`
	RegisteredAt float64      `json:"registered_at" db:"registered_at"`
	CheckImage   []CheckImage `json:"check_image,omitempty"`
	CheckList    []string
}

type CheckImage struct {
	Id         int    `json:"id" db:"id"`
	CheckId    int    `json:"check_id" db:"check_id"`
	CheckImage string `json:"check_image" db:"check_image"`
}

type GetAllChecksCategoryResponse struct {
	Data     []*UserChecks  `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}

type FilterForCheck struct {
	PartnerId           int     `json:"partner_id" query:"partner_id" form:"partner_id"`
	MoneyAmount         int     `json:"money_amount" query:"money_amount" form:"money_amount"`
	StartCheckDate      float64 `json:"start_check_date" query:"start_check_date" form:"start_check_date"`
	EndCheckDate        float64 `json:"end_check_date" query:"end_check_date" form:"end_check_date"`
	StartRegisteredDate float64 `json:"start_register_date" query:"start_register_date" form:"start_register_date"`
	EndRegisteredDate   float64 `json:"end_register_date" query:"end_register_date" form:"end_register_date"`
}
