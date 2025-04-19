package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	converter "github.com/Ranik23/avito-tech-spring/internal/models/converter/http"
	"github.com/Ranik23/avito-tech-spring/internal/models/dto"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type PvzController interface {
	CreatePvz(c *gin.Context)
	GetPvzInfo(c *gin.Context)
	GetPvzInfoOptimized(c *gin.Context)
	CloseLastReception(c *gin.Context)
	DeleteLastProduct(c *gin.Context)
	CreateReception(c *gin.Context)
	AddProduct(c *gin.Context)
}

type pvzController struct {
	service service.Service
	logger *slog.Logger
}

func NewPVZController(service service.Service) PvzController {
	return &pvzController{
		service: service,
	}
}

func (p *pvzController) CloseLastReception(c *gin.Context) {
	if !p.check(c, "employee") {
		return
	}

	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: "no pvzId provided",
		})
		return
	}

	reception, err := p.service.CloseReception(c, pvzID)
	if err != nil {
		if errors.Is(err, service.ErrAllReceptionsClosed) {
			c.JSON(http.StatusBadGateway, dto.Error{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, converter.FromDomainReceptionToCloseReseptionResp(reception))
}


func (p *pvzController) AddProduct(c *gin.Context) {
	if !p.check(c, "employee") {
		return
	}

	var req dto.PostProductReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	product, err := p.service.AddProduct(c, req.PvzID, req.Type)
	if err != nil {
		if errors.Is(err, service.ErrAllReceptionsClosed) {
			c.JSON(http.StatusBadRequest, dto.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}
	converter.FromDomainProductToDtoPostProductResp(product)

	c.JSON(http.StatusCreated, converter.FromDomainProductToDtoPostProductResp(product))
}

func (p *pvzController) CreateReception(c *gin.Context) {
	if !p.check(c, "employee") {
		return
	}

	var req dto.CreateReceptionReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	reception, err := p.service.StartReception(context.TODO(), req.PvzId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.FromDomainReceptionToCreateReceptionResp(reception))
}


func (p *pvzController) DeleteLastProduct(c *gin.Context) {
	if !p.check(c, "employee") {
		return
	}

	pvzID := c.Param("pvzId")
	if pvzID == "" {
		c.JSON(http.StatusBadRequest,dto.Error{Message: "pvzID not provided"})
		return
	}

	if err := p.service.DeleteLastProduct(c, pvzID); err != nil {
		if errors.Is(err, service.ErrAllReceptionsClosed) || errors.Is(err, service.ErrReceptionEmpty) {
			c.JSON(http.StatusBadRequest, dto.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Description" : "Товар Удален",
	})
}


func (p *pvzController) CreatePvz(c *gin.Context) {
	
	if !p.check(c, "moderator") {
		return
	}

	var req dto.CreatePvzReq 

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	pvz, err := p.service.CreatePVZ(c, req.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, converter.FromDomainPVZToCreatePvzResp(pvz))
}

func (p *pvzController) GetPvzInfoOptimized(c *gin.Context) {
	if !p.check(c, "employee", "moderator") {
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 30 {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
		return
	}

	startDate, ok := p.parseTimeParam(c, "startDate")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: "failed to get startDate",
		})
		return
	}

	endDate, ok := p.parseTimeParam(c, "endDate")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: "failed to get endDate",
		})
		return
	}

	_, err = p.service.GetPVZSInfoOptimized(c, startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{
			Message: err.Error(),
		})
		return
	}
}

func (p *pvzController) GetPvzInfo(c *gin.Context) {
	if !p.check(c, "employee", "moderator") {
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 30 {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: err.Error(),
		})
		return
	}

	startDate, ok := p.parseTimeParam(c, "startDate")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: "failed to get startDate",
		})
		return
	}

	endDate, ok := p.parseTimeParam(c, "endDate")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.Error{
			Message: "failed to get endDate",
		})
		return
	}

	_, err = p.service.GetPVZSInfo(c, startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{
			Message: err.Error(),
		})
		return
	}
}


func (p *pvzController) parseTimeParam(c *gin.Context, param string) (time.Time, bool) {
	value := c.Query(param)
	if value == "" {
		return time.Now(), true
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Message: fmt.Sprintf("invalid %s format", param)})
		return time.Now(), false
	}
	return t, true
}


func (p *pvzController) check(c *gin.Context, allowedRoles ...string) bool {
	role, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
			Message: "no role provided",
		})
		return false
	}

	roleStr, ok := role.(string) 
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.Error{
			Message: "role must be a string",
		})
		return false
	}

	roleStr = strings.ToLower(roleStr)

	for _, allowedRole := range allowedRoles {
		if roleStr == strings.ToLower(allowedRole) {
			return true
		}
	}

	c.AbortWithStatusJSON(http.StatusForbidden, dto.Error{
		Message: "this role is forbidden",
	})

	return false
}

