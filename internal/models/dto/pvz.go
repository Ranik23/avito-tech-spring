package dto


type Pvz struct {
	Id				 string		`json:"id"`
	City 			 string		`json:"city"`
	RegistrationDate string		`json:"registrationDate"`
}


type GetPvzResp struct {
	Id				 string		`json:"id"`
	City 			 string		`json:"city"`
	RegistrationDate string		`json:"registrationDate"`	
}


type CreatePvzReq struct {
	Id				 string		`json:"id"`
	City 			 string		`json:"city"`
	RegistrationDate string		`json:"registrationDate"`
}


type CreatePvzResp struct {
	Id				 string		`json:"id"`
	City 			 string		`json:"city"`
	RegistrationDate string		`json:"registrationDate"`
}