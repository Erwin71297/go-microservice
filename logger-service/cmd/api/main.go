package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	corsConfig.MaxAge = 300
	return corsConfig
}

type Config struct {
	Models data.Models
}

func main() {
	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	//Setup gin
	app := gin.Default()
	app.Use(cors.New(CORSConfig()))
	Routes(app)

	//Register the RPC Server
	err := rpc.Register(new(RPCServer))
	go RPCListen()

	go gRPCListen()

	//Setup Listen and serve
	log.Println("Starting service on port: ", webPort)
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app,
	}

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func Serve(c *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: c,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func RPCListen() error {
	log.Println("Starting RPC server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Println("error rpcConn", err)
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Connecting: ", err)
		return nil, err
	}
	log.Println("Connected to Mongo")

	return conn, nil
}
