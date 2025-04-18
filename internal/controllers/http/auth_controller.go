package http

import (
	"context"
	"errors"

	"github.com/Ranik23/avito-tech-spring/internal/models/dto"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	DummyLogin(c *gin.Context)
}

type authController struct {
	service service.Service
}


func NewAuthController(service service.Service) AuthController {
	return &authController{
		service: service,
	}
}


func (a *authController) DummyLogin(c *gin.Context) {
	
	var req dto.DummyLoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}


	token, err := a.service.DummyLogin(context.TODO(), req.Role)
	if err != nil {
		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.DummyLoginResp{
		Token: token,
	})
}

func (a *authController) Login(c *gin.Context) {

	var req dto.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}

	token, err := a.service.Login(context.TODO(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(401, dto.Error{
				Message: err.Error(),
			})
			return
		}

		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.LoginResp{
		Token: token,
	})
}

func (a *authController) Register(c *gin.Context) {

	var req dto.RegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}


	userID, err := a.service.Register(context.TODO(), req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(201, dto.RegisterResp{
		Email: req.Email,
		Id: userID,
		Role: req.Role,
	})
}
