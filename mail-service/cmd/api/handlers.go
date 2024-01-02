package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})

}

func SendMail(c *gin.Context) {
	Mailer := CreateMail()

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := ReadJSON(c, &requestPayload)
	if err != nil {
		log.Println("error read json")
		ErrorJSON(c, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println("send smtp error", err)
		ErrorJSON(c, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Send to" + requestPayload.To,
	}

	c.JSON(http.StatusAccepted, payload)
}
