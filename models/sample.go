package models

type SampleRes struct {
	Id           int32  `json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	Number       string `json:"number" form:"number"`
	CountryId    int    `json:"country_id" form:"country_id"`
	Code         int32  `json:"code" form:"code"`
	TermAccepted bool   `json:"term_accepted" form:"term_accepted"`
}

type SamplePostReq struct {
	Id           int32  `json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	Number       string `json:"number" form:"number"`
	CountryId    int    `json:"country_id" form:"country_id"`
	Code         int32  `json:"code" form:"code"`
	TermAccepted bool   `json:"term_accepted" form:"term_accepted"`
}
type SamplePostRes struct {
	Data SamplePostReq `json:"data" form:"data"`
}

type SamplePatchReq struct {
	Id           int32  `json:"id" form:"id"`
	Name         string `json:"name" form:"name"`
	Number       string `json:"number" form:"number"`
	CountryId    int    `json:"country_id" form:"country_id"`
	Code         int32  `json:"code" form:"code"`
	TermAccepted bool   `json:"term_accepted" form:"term_accepted"`
}
