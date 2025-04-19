package app

import (
	"github.com/Ranik23/avito-tech-spring/internal/controllers/http"
	"github.com/Ranik23/avito-tech-spring/internal/controllers/http/middleware"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/gin-gonic/gin"
)







func SetUpRoutes(router *gin.Engine, authController http.AuthController, pvzController http.PvzController, tokenService token.Token) {

	router.POST("/dummyLogin", authController.DummyLogin)
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	group := router.Group("/")
	group.Use(middleware.JwtAuth(tokenService))
	{
		group.POST("/pvz", pvzController.CreatePvz)
		group.POST("/receptions", pvzController.CreateReception)
		group.POST("/products", pvzController.AddProduct)
		group.POST("/pvz/:pvzId/delete_last_product", pvzController.DeleteLastProduct)
		group.POST("/pvz/:pvzId/close_last_reception", pvzController.CloseLastReception)
		group.GET("/pvz", pvzController.GetPvzInfo)
	}
}