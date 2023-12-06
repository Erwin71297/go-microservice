package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
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
	var w http.ResponseWriter
	//var r *http.Request

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = WriteJSON(w, http.StatusOK, payload)
}

func HandleSubmission(c *gin.Context) {
	var w http.ResponseWriter
	var r *http.Request
	var requestPayload RequestPayload

	err := ReadJSON(w, r, &requestPayload)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		Authenticate(c, w, requestPayload.Auth)
	default:
		ErrorJSON(w, errors.New("unknown action"))
	}
}

func Authenticate(c *gin.Context, w http.ResponseWriter, a AuthPayload) {

	//Create some json that will be sent to the auth microservices
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//Call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the correct status
	if response.StatusCode == http.StatusUnauthorized {
		ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	//Create variable to read response.Body
	var jsonFromService jsonResponse

	//Decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	WriteJSON(w, http.StatusAccepted, payload)
}

func Testing(c *gin.Context) {
	//var r *http.Request

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	if err := c.BindJSON(&payload); err != nil {
		log.Panic("error:", err)
		return
	}

	c.JSON(http.StatusAccepted, payload)
}

func Ping(c *gin.Context) {
	log.Println("pings")

	c.JSON(200, gin.H{
		"message": "pong",
	})

}
