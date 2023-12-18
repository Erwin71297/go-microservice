package main

import (
	"authentication/data"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	log.Println("enter authenticate-service")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	requestPayload.Email = "admin@example.com"
	requestPayload.Password = "verysecret"
	log.Println("request Payload", requestPayload)

	err := ReadJSON(c, &requestPayload)
	if err != nil {
		log.Println("it enters here", err)
		ErrorJSON(c, err, http.StatusBadRequest)
		return
	}

	//Validate the user against the database
	db := data.New(connectToDB())
	user, err := db.User.GetByEmail(c, requestPayload.Email)
	if err != nil {
		log.Println("error email ", err)
		ErrorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Println("error password ", err)
		ErrorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	//WriteJSON(c, http.StatusAccepted, payload)
	c.JSON(http.StatusAccepted, payload)
}

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
