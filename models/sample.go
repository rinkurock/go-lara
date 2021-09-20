package models

type SamplePostReq struct {
	Name         string `json:"name" form:"name"`
	Number       string `json:"number" form:"number"`
	CountryId    int    `json:"country_id" form:"country_id"`
	Code         int32  `json:"code" form:"code"`
	TermAccepted bool   `json:"term_accepted" form:"term_accepted"`
}

type SamplePostRes struct {
	Data SamplePostReq `json:"data" form:"data"`
}
