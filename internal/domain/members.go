package domain

type Winners struct {
	Id          int     `json:"id" db:"id"`
	CheckId     int     `json:"check_id,omitempty" db:"check_id"`
	PhoneNumber string  `json:"phone_number" db:"phone_number"`
	RaffleType  string  `json:"raffle_type,omitempty" db:"raffle_type"`
	RaffleDate  float64 `json:"raffle_date" db:"raffle_date"`
	UserName    *string `json:"user_name" db:"user_name"`
	IsWinner    bool    `json:"is_winner" db:"is_winner"`
}

type DayInput struct {
	CreatedAt float64 `json:"created_at" db:"created_at"`
}

type GetAllDaysResponse struct {
	Data     []*DayInput    `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}

type Members struct {
	Id       int  `json:"id" db:"id"`
	CheckId  int  `json:"check_id" db:"check_id"`
	RaffleId int  `json:"raffle_id" db:"raffle_id"`
	IsWinner bool `json:"is_winner" db:"is_winner"`
}
