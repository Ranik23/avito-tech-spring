package dto



type DummyLoginReq struct {
	Role string `json:"role"`
}


type DummyLoginResp struct {
	Token string		`json:"token"`
}