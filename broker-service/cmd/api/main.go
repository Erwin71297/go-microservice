package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const webPort = "80"

type Config struct{}

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

	//Specify cors
	app.Use(cors.New(CORSConfig()))
	GinRoutes(app)

	log.Printf("Starting broker service on port %s\n", webPort)

	//Define routing

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

}
