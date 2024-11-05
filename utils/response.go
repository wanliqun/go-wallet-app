package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
		"result":  data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}
