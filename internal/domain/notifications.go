package domain

type Notification struct {
	Id        int           `json:"id" db:"id"`
	Title     string        `json:"title" db:"title"`
	Text      string        `json:"text" db:"text" `
	PartnerId int           `json:"partner_id,omitempty" db:"partner_id"`
	Logo      string        `json:"logo,omitempty" db:"logo"`
	Link      string        `json:"link,omitempty" db:"link"`
	Reference string        `json:"reference,omitempty" db:"reference"`
	Date      float64       `json:"date" db:"noty_date"`
	Time      int           `json:"time" db:"noty_time"`
	Status    int           `json:"status,omitempty" db:"status"`
	Getters   int           `json:"getters,omitempty" db:"noty_getters"`
	Ids       []int         `json:"ids,omitempty"`
	Users     []*GetterList `json:"users" db:"users"`
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
	UserId            *int   `json:"user_id" db:"user_id"`
}
