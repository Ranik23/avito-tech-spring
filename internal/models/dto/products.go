package dto



type Product struct {
	DateTime    string `json:"dateTime,omitempty"`  
	Id          string `json:"id,omitempty"`     
	ReceptionID string `json:"receptionId"`
	Type        string `json:"type"`
}


type PostProductReq struct {
	PvzID 	string		`json:"pvzId"`
	Type 	string		`json:"type"`	
}

type PostProductResp struct {
	DateTime    string `json:"dateTime,omitempty"`  
	Id          string `json:"id,omitempty"`     
	ReceptionID string `json:"receptionId"`
	Type        string `json:"type"`
}
