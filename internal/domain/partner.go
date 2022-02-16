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

type UpdatePartner struct {
	PartnerName      string  `json:"partner_name" form:"partner_name"`
	Position         int     `json:"position" form:"position"`
	Logo             *string `json:"logo" form:"logo"`
	LinkWebsite      string  `json:"link_website" form:"link_website"`
	Banner           *string `json:"banner" form:"banner"`
	BannerKz         *string `json:"banner_kz" form:"banner_kz"`
	Status           int     `form:"status"  enums:"1,2" default:"1"`
	StartPartnership string  `json:"start_partnership" form:"start_partnership"`
	EndPartnership   string  `json:"end_partnership" form:"end_partnership"`
	PartnerPackage   string  `json:"partner_package" form:"partner_package"`
	Reference        string  `json:"reference" form:"reference"`
}

type GetAllPartnersCategoryResponse struct {
	Data     []*Partner     `json:"data"`
	PageInfo PaginationPage `json:"page-info"`
}
