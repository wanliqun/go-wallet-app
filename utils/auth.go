package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")

func ExtractBearerToken(c *gin.Context) (string, error) {
	token := strings.TrimSpace(c.GetHeader("Authorization")) // Trim any extra spaces

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(token, bearerPrefix) {
		return "", ErrInvalidAuthorizationHeader
	}

	return strings.TrimPrefix(token, bearerPrefix), nil
}
