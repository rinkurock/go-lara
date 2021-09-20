package repo

import m "app/models"

func GetPostResponse(req m.SamplePostReq) m.SamplePostRes {
	data := m.SamplePostRes{}

	data.Data = req

	return data

}
