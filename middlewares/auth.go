package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wanliqun/go-wallet-app/services"
	"github.com/wanliqun/go-wallet-app/utils"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, services.ErrUnauthorized)
			return
		}

		user, ok, err := userService.GetUserByName(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, services.ErrUserNotFound)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
