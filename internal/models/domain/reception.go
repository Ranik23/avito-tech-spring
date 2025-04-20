package domain

import "time"

type Reception struct {
	ID       string
	DateTime time.Time
	PvzID    string
	Status   string
}
