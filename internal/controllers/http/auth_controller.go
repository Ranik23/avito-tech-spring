package http

import (
	"errors"
	"log/slog"

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
	logger  *slog.Logger
}

func NewAuthController(service service.Service, logger *slog.Logger) AuthController {
	return &authController{
		service: service,
		logger:  logger,
	}
}

func (a *authController) DummyLogin(c *gin.Context) {
	var req dto.DummyLoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("DummyLogin failed", slog.String("error", err.Error()))
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("DummyLogin request received", slog.String("role", req.Role))

	token, err := a.service.DummyLogin(c, req.Role)
	if err != nil {
		a.logger.Error("DummyLogin service error", slog.String("error", err.Error()))
		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("DummyLogin successful", slog.String("token", token))
	c.JSON(200, dto.DummyLoginResp{
		Token: token,
	})
}

func (a *authController) Login(c *gin.Context) {
	var req dto.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("Login failed", slog.String("error", err.Error()))
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("Login request received", slog.String("email", req.Email))

	token, err := a.service.Login(c, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			a.logger.Warn("Invalid credentials", slog.String("email", req.Email))
			c.JSON(401, dto.Error{
				Message: err.Error(),
			})
			return
		}

		a.logger.Error("Login service error", slog.String("error", err.Error()))
		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("Login successful", slog.String("email", req.Email), slog.String("token", token))
	c.JSON(200, dto.LoginResp{
		Token: token,
	})
}

func (a *authController) Register(c *gin.Context) {
	var req dto.RegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("Register failed", slog.String("error", err.Error()))
		c.JSON(400, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("Register request received", slog.String("email", req.Email), slog.String("role", req.Role))

	userID, err := a.service.Register(c, req.Email, req.Password, req.Role)
	if err != nil {
		a.logger.Error("Register service error", slog.String("error", err.Error()))
		c.JSON(500, dto.Error{
			Message: err.Error(),
		})
		return
	}

	a.logger.Info("Register successful", slog.String("email", req.Email), slog.String("userID", userID))
	c.JSON(201, dto.RegisterResp{
		Email: req.Email,
		Id:    userID,
		Role:  req.Role,
	})
}
