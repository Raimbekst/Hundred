package domain

type User struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"user_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
	Age         int    `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	UserType    string `json:"user_type,omitempty" db:"user_type"`
	Password    string `json:"password,omitempty" db:"password"`
	CityId      int    `json:"city" db:"city_id"`
	IsBlocked   bool   `json:"is_blocked" db:"is_blocked"`
}

type UserList struct {
	Id           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"user_name"`
	PhoneNumber  string  `json:"phone_number" db:"phone_number"`
	Email        string  `json:"email" db:"email"`
	Age          int     `json:"age" db:"age"`
	Gender       string  `json:"gender" db:"gender"`
	UserType     string  `json:"user_type,omitempty" db:"user_type"`
	City         *string `json:"city" db:"city"`
	IsBlocked    bool    `json:"is_blocked" db:"is_blocked"`
	RegisteredAt float64 `json:"registered_at" db:"registered_at"`
}

type GetAllUsersCategoryResponse struct {
	Data     []*UserList    `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}

type EmailInput struct {
	Email string `json:"email" validate:"required"`
}

type UserChecks struct {
	Id           int           `json:"id" db:"id"`
	UserId       int           `json:"user_id" db:"user_id"`
	UserName     string        `json:"user_name" db:"user_name"`
	PhoneNumber  string        `json:"phone_number" db:"phone_number"`
	IsBlocked    bool          `json:"is_blocked" db:"is_blocked"`
	PartnerName  string        `json:"partner_name" db:"partner_name"`
	CheckAmount  string        `json:"check_amount" db:"check_amount"`
	CheckDate    float64       `json:"check_date" db:"check_date"`
	CheckImage   []*CheckImage `json:"check_image" db:"checkImg"`
	RegisteredAt float64       `json:"registered_at" db:"registered_at"`
}
