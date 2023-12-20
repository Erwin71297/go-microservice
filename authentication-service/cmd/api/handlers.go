package main

import (
	"authentication/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var requestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Authenticate(c *gin.Context) {
	log.Println("enter authenticate-service")

	if requestPayload.Email == "" {
		requestPayload.Email = "admin@example.com"
		requestPayload.Password = "verysecret"
	}
	log.Println("request Payload", requestPayload)

	_ = ReadJSON(c, &requestPayload)

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

	// Log authentication
	err = LogRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		ErrorJSON(c, err)
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

func LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = "Name"
	entry.Data = "Data"

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
