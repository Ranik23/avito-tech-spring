package dto



type RegisterReq struct {
	Email 		string		`json:"email"`
	Password 	string		`json:"password"`
	Role		string		`json:"role"`
}


type RegisterResp struct {
	Email 	string		`json:"email"`
	Id		string		`json:"id,omitempty"`
	Role	string		`json:"role"`
}