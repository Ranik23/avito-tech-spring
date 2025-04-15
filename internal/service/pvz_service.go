package service


type PVZService interface {

}


type pvzService struct {

}

func NewPVZService() PVZService {
	return &pvzService{}
}