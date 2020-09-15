package routes

import (
	"sprout_server/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// register the routes

	//eg.
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}
