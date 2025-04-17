package service

import "errors"


var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound = errors.New("not found")
	ErrAlreadyOpen = errors.New("already open")
	ErrAllReceptionsClosed = errors.New("all receptions closed")
	ErrEmpty = errors.New("empty")
	ErrNoPVZFound = errors.New("no pvz found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRole = errors.New("invalid role")
)



type Service interface {
	AuthService
	PVZService
}