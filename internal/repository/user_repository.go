package repository

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context)
	GetUsers(ctx context.Context)
}

type userRepository struct {
}

// GetUsers implements UserRepository.
func (u *userRepository) GetUsers(ctx context.Context) {
	panic("unimplemented")
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(ctx context.Context) {
	panic("unimplemented")
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
