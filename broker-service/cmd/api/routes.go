package main

import (
	"github.com/gin-gonic/gin"
)

func GinRoutes(app *gin.Engine) {
	//Heartbeat
	app.GET("/ping", Ping)

	//Routes
	app.POST("/broker", Broker)
	app.POST("/handler", HandleSubmission)
}
