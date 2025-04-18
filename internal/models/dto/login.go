package dto



type LoginReq struct {
	Email string		`json:"email"`
	Password string		`json:"password"`
}

type LoginResp struct {
	Token string		`yaml:"token"`
}

