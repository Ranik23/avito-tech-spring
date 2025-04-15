package repository

import "time"


// Сущности репозитория

type Product struct {
	ID 		string
	AddDate time.Time
	Type 	string
}

type PVZ struct {
	ID 			uint
	City 		string
	RegDate 	time.Time
}

type Reception struct {
	ID 			string
	Status  	string
	PvzID		string
	StartDate 	time.Time
	CloseDate 	time.Time
}

type User struct {
	ID           string 
	Email        string 
	PasswordHash string
	Role         string 
	RegDate    	 time.Time 
}
