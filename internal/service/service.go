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
	ErrInvalidCity = errors.New("invalid city")
	ErrUserNotFound = errors.New("user not found")
	ErrReceptionEmpty = errors.New("reception is empty")
)

type Service interface {
	AuthService
	PVZService
}

type service struct {
	AuthService
	PVZService
}


func NewService(authService AuthService, pvzService PVZService) Service {
	return &service{
		AuthService: authService,
		PVZService: pvzService,
	}
}