package http

import (
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type PvzController interface {
	CreatePvz(c *gin.Context)
	GetPvz(c *gin.Context)
	CloseLastReception(c *gin.Context)
	DeleteLastProduct(c *gin.Context)
	CreateReception(c *gin.Context)
	CreateProduct(c *gin.Context)
}

type pvzController struct {
	srv service.Service
}

func NewPVZController(srv service.Service) PvzController {
	return &pvzController{
		srv: srv,
	}
}

// CloseLastReception implements PvzController.
func (p *pvzController) CloseLastReception(c *gin.Context) {
	panic("unimplemented")
}

// CreateProduct implements PvzController.
func (p *pvzController) CreateProduct(c *gin.Context) {
	panic("unimplemented")
}

// CreateReception implements PvzController.
func (p *pvzController) CreateReception(c *gin.Context) {
	panic("unimplemented")
}

// DeleteLastProduct implements PvzController.
func (p *pvzController) DeleteLastProduct(c *gin.Context) {
	panic("unimplemented")
}

// CreatePvz implements PvzController.
func (p *pvzController) CreatePvz(c *gin.Context) {
	panic("unimplemented")
}

// GetPvz implements PvzController.
func (p *pvzController) GetPvz(c *gin.Context) {
	panic("unimplemented")
}
