package domain

type Partner struct {
	Id               int    `json:"id" db:"id"`
	Position         int    `json:"position" db:"position"`
	PartnerName      string `json:"name" db:"partner_name"`
	Logo             string `json:"logo" db:"logo"`
	LinkWebsite      string `json:"linkWebsite" db:"link_website"`
	Banner           string `json:"banner" db:"banner"`
	BannerKz         string `json:"banner_kz" db:"banner_kz"`
	Status           int    `json:"status" db:"status"`
	StartPartnership string `json:"start_partnership" db:"start_partnership"`
	EndPartnership   string `json:"end_partnership" db:"end_partnership"`
	PartnerPackage   string `json:"partner_package" db:"partner_package"`
	Reference        string `json:"reference" db:"reference"`
}

type GetAllPartnersCategoryResponse struct {
	Data     []*Partner     `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
