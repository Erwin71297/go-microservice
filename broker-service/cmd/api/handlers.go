package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Broker(c *gin.Context) {
	var response jsonResponse

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	response = payload
	response.Data = "Broker was hit"

	c.Bind(response)
	c.JSON(200, response)
}

func HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	var newPayload RequestPayload
	newPayload.Action = "auth"
	newPayload.Auth.Email = "admin@example.com"
	newPayload.Auth.Password = "verysecret!"

	if requestPayload.Action == "" {
		requestPayload = newPayload
	}

	err := ReadJSON(c, &requestPayload)
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		Authenticate(c, requestPayload.Auth)
	default:
		ErrorJSON(c, errors.New("unknown action"))
	}
}

func Authenticate(c *gin.Context, a AuthPayload) {

	//Create some json that will be sent to the auth microservices
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//Call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the correct status
	if response.StatusCode == http.StatusUnauthorized {
		ErrorJSON(c, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ErrorJSON(c, errors.New("error calling auth service"))
		return
	}

	//Create variable to read response.Body
	var jsonFromService jsonResponse

	//Decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	if jsonFromService.Error {
		ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	// WriteJSON(c, http.StatusAccepted, payload)
	c.Bind(payload)
	c.JSON(200, payload)
}

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})

}
