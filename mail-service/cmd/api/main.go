package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const webPort = "80"

type Config struct {
	Mailer Mail
}

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
	//Setup gin
	app := gin.Default()
	app.Use(cors.New(CORSConfig()))
	// app.Use(ApiMiddleware(Config{
	// 	Mailer: createMail(),
	// }))
	Routes(app)

	log.Println("Starting mail service on port", webPort)

	//Setup Listen and serve
	//go Serve(app)
	log.Println("Starting service on port: ", webPort)
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func CreateMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_Name"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
