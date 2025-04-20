package middleware

import (
	"net/http"

	"github.com/Ranik23/avito-tech-spring/internal/models/dto"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/gin-gonic/gin"
)

func JwtAuth(tokenService token.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
				Message: "no token provided",
			})
			return
		}
		
		claims, err := tokenService.Parse(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
				Message: err.Error(),
			})
			return
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}
