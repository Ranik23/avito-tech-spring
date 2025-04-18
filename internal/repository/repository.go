package repository

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	UserRepository
	PvzRepository
	ProductRepository
}

type repository struct {
	UserRepository
	PvzRepository
	ProductRepository
}

func NewRepository(userRepo UserRepository, pvzRepo PvzRepository,
				productRepo ProductRepository, receptionRepo ReceptionRepository) Repository {
	return &repository{
		UserRepository: userRepo,
		PvzRepository: pvzRepo,
		ProductRepository: productRepo,
	}
}