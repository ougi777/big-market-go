package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
