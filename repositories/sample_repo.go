package repo

import m "app/models"

func getSampleData() m.SampleRes {
	d := m.SampleRes{}
	d.Id = 99
	d.Name = "John Doe"
	d.CountryId = 1
	d.Number = "+880111111111111"
	d.TermAccepted = true
	d.Code = 100000
	return d
}

func SampleGetResponse() m.SamplePostRes {
	d := m.SamplePostRes{}
	data := getSampleData()
	d.Data.Id = data.Id
	d.Data.Name = data.Name
	d.Data.Number = data.Number
	d.Data.CountryId = data.CountryId
	d.Data.Code = data.Code
	d.Data.TermAccepted = data.TermAccepted
	return d
}

func GetPostResponse(req m.SamplePostReq) m.SamplePostRes {
	data := m.SamplePostRes{}
	data.Data = req
	data.Data.Id = 99
	return data
}
func PatchPesponse(req m.SamplePatchReq) m.SamplePostRes {
	d := m.SamplePostRes{}
	data := getSampleData()
	d.Data.Id = data.Id
	d.Data.Name = data.Name
	d.Data.Number = data.Number
	d.Data.CountryId = data.CountryId
	d.Data.Code = data.Code
	d.Data.TermAccepted = data.TermAccepted

	if req.Name != data.Name {
		d.Data.Name = req.Name
	}
	if req.Number != data.Number {
		d.Data.Number = req.Number
	}
	if req.CountryId != data.CountryId {
		d.Data.CountryId = req.CountryId
	}

	if req.Code != data.Code {
		d.Data.Code = req.Code
	}
	if req.TermAccepted != data.TermAccepted {
		d.Data.TermAccepted = req.TermAccepted
	}
	return d
}
