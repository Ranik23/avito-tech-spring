package http

import (
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	DummyLogin(c *gin.Context)
}

type authController struct {
	srv service.Service
}


func NewAuthController(srv service.Service) AuthController {
	return &authController{
		srv: srv,
	}
}


func (a *authController) DummyLogin(c *gin.Context) {
	panic("unimplemented")
}

func (a *authController) Login(c *gin.Context) {
	panic("unimplemented")
}

func (a *authController) Register(c *gin.Context) {
	panic("unimplemented")
}
