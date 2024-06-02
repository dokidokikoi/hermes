package ggbases

type DetailResp struct {
	Category   string `json:"category"`
	CreateDate int64  `json:"create_date"`
	UpdateDate int64  `json:"update_date"`
	HomeStatus string `json:"home_status"`
	HomeUrl    string `json:"home_url"`
}
