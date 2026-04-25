package router

import (
	"transport/api/src/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	targetServiceHandler *handler.TargetServiceHandler,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		targetServices := api.Group("/target-services")
		targetServices.POST("/", targetServiceHandler.Create)
		targetServices.GET("/:id", targetServiceHandler.Get)
		targetServices.DELETE("/:id", targetServiceHandler.Delete)
	}

	return router
}
