package service

import "time"


type PVZ struct {
	ID		string		`json:"id"`
	RegDate	time.Time	`json:"reg_date"`
	City	string		`json:"city"`
}

type PVZInfo struct {
	Pvz PVZ
	Receptions []ReceptionInfo
}

type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    string    `json:"pvz_id"`
	Status   string    `json:"status"`
}

type ReceptionInfo struct {
	Reception Reception
	Products  []Product
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"date_time"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"reception_id"`
}