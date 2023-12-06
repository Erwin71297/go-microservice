package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	router "github.com/go-microservice/broker-service/cmd/api/routes"
)

const webPort = "80"

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	corsConfig.MaxAge = 300
	return corsConfig
}

func main() {
	gin.SetMode(gin.DebugMode)
	app := gin.Default()
	app.Use(gin.Logger())

	//Specify cors
	app.Use(cors.New(CORSConfig()))

	log.Printf("Starting broker service on port %s\n", webPort)

	//Define routing
	app.GET("/ping")
	router.GinRoutes(&app.RouterGroup)

}
