package repository


type Repository interface {
	UserRepository
}

type repository struct {
	UserRepository
}


func NewRepository() Repository {
	return &repository{}
}