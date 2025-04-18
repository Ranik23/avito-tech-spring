package dto


type Reception struct {
	DateTime	string		`json:"dateTime"`
	Id			string		`json:"id,omitempty"`
	PvzId		string		`json:"pvzId"`
	Status		string		`json:"status"`
}


type CloseReceptionResp struct {
	DateTime	string		`json:"dateTime"`
	Id			string		`json:"id,omitempty"`
	PvzId		string		`json:"pvzId"`
	Status		string		`json:"status"`
}

type CreateReceptionReq struct {
	PvzId string		`json:"pvzId"`
}


type CreateReceptionResp struct {
	DateTime	string		`json:"dateTime"`
	Id			string		`json:"id,omitempty"`
	PvzId		string		`json:"pvzId"`
	Status		string		`json:"status"`
}