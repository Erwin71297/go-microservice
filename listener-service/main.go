package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//Connect to rabbitMQ
	rabbitConn, err := Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//Listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	//Create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	//Watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println("error Listening", err)
	}
}

func Connect() (*amqp.Connection, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

	// Don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
			os.Exit(0)
		} else {
			log.Println("Connected to rabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 6)) * time.Second
		log.Println("backing off... ")
		time.Sleep(backOff)
	}

	return connection, nil
}
