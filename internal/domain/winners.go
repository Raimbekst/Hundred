package domain

type WinnerInput struct {
	CheckId             int     `json:"check_id"    validate:"required"`
	RaffleId            int     `json:"raffle_id"   validate:"required"`
	PartnerId           int     `json:"partner_id" db:"partner_id"`
	MoneyAmount         int     `json:"money_amount" db:"check_amount"`
	StartCheckDate      float64 `json:"start_check_date"`
	EndCheckDate        float64 `json:"end_check_date"`
	StartRegisteredDate float64 `json:"start_register_date"`
	EndRegisteredDate   float64 `json:"end_register_date"`
}

type GetAllWinnersCategoryResponse struct {
	Data     []*Winners     `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
