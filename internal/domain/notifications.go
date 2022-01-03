package domain

type Notification struct {
	Id        int         `json:"id" db:"id"`
	Title     string      `json:"title" db:"title" validate:"required"`
	Text      string      `json:"text" db:"text" validate:"required"`
	PartnerId int         `json:"partner_id" db:"partner_id"`
	Link      string      `json:"link" db:"link"`
	Reference string      `json:"reference" db:"reference"`
	Date      float64     `json:"date" db:"noty_date"`
	Time      float32     `json:"time" db:"noty_time"`
	Status    string      `json:"status,omitempty" db:"status"`
	Getters   string      `json:"getters,omitempty" db:"noty_getters"`
	Ids       []int       `json:"ids,omitempty"`
	Users     *GetterList `json:"users" db:"users"`
}

type GetterList struct {
	Id             int `json:"id" db:"id"`
	NotificationId int `json:"notification_id" db:"notification_id"`
	UserId         int `json:"user_id" db:"user_id"`
}

type GetAllNotificationsResponse struct {
	Data     []*Notification `json:"data"`
	PageInfo PaginationPage  `json:"page-info"`
}

type NotificationToken struct {
	Id                int    `json:"id" db:"id"`
	RegistrationToken string `json:"registration_token" db:"registration_token"`
	UserId            int    `json:"user_id" db:"user_id"`
}
