package domain

type Raffle struct {
	Id            int     `json:"id" db:"id"`
	UserId        *int    `json:"user_id" db:"user_id"`
	CheckId       *int    `json:"check_id" db:"check_id"`
	RaffleDate    float64 `json:"raffle_date" db:"raffle_date"`
	RaffleTime    int     `json:"raffle_time" db:"raffle_time"`
	CheckCategory int     `json:"check_category" db:"check_category"`
	RaffleType    int     `json:"raffle_type" db:"raffle_type"`
	Status        string  `json:"status" db:"status"`
	Reference     string  `json:"reference"  db:"reference"`
	UserName      *string `json:"user_name" db:"user_name" default:""`
	PhoneNumber   *string `json:"phone_number" db:"phone_number" default:""`
}

type GetAllRaffleCategoryResponse struct {
	Data     []*Raffle      `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}

type FilterForRaffles struct {
	IsFinished int     `json:"is_finished" query:"is_finished" form:"is_finished" enums:"1,2"`
	RaffleType int     `json:"raffle_type" query:"raffle_type" form:"raffle_type" enums:"1,2,3" example:"1"`
	RaffleDate float64 `json:"raffle_date" query:"raffle_date" form:"raffle_date"`
	RaffleTime int     `json:"raffle_time" query:"raffle_time" form:"raffle_time"`
}
