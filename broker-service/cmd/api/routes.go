package main

import (
	"github.com/gin-gonic/gin"
)

func GinRoutes(app *gin.RouterGroup) {
	api := app.Group("/v1")

	//Heartbeat
	api.GET("/ping")

	//Routes
	api.POST("", Broker)
	api.POST("/test", Testing)
	api.POST("/handle", HandleSubmission)

	return
}
