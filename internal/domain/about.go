package domain

type AboutUs struct {
	Id            int    `json:"-" db:"id"`
	FacebookLink  string `json:"facebook_link" db:"facebook_link"`
	YoutubeLink   string `json:"youtube_link" db:"youtube_link"`
	InstagramLink string `json:"instagram_link" db:"instagram_link"`
	TiktokLink    string `json:"tiktok_link" db:"tiktok_link"`
	WhatsappLink  string `json:"whatsapp_link" db:"whatsapp_link"`
	TelegramLink  string `json:"telegram_link" db:"telegram_link"`
	PhoneNumber   string `json:"phone_number" db:"phone_number"`
	PhoneNumber2  string `json:"phone_number_2" db:"phone_number_2"`
}

type GetAllAboutUsCategoryResponse struct {
	Data     []*AboutUs     `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
