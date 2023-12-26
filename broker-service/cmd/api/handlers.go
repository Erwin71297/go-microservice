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
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func Broker(c *gin.Context) {
	var response jsonResponse

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	response = payload
	response.Data = "Broker was hit"

	c.JSON(200, response)
}

func HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	err := ReadJSON(c, &requestPayload)
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		log.Println("enter here")
		Authenticate(c, requestPayload.Auth)
	case "log":
		LogItem(c, requestPayload.Log)
	case "mail":
		SendMail(c, requestPayload.Mail)
	default:
		ErrorJSON(c, errors.New("unknown action"))
	}
}

func LogItem(c *gin.Context, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		ErrorJSON(c, errors.New("error log service"+err.Error()))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	// WriteJSON(c, http.StatusAccepted, payload)
	c.JSON(http.StatusAccepted, payload)
}

func Authenticate(c *gin.Context, a AuthPayload) {
	//Create some json that will be sent to the auth microservices
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	authServiceUrl := "http://authentication-service/authenticate"

	//Call the service
	request, err := http.NewRequest("POST", authServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the correct status
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("error code 1")
		ErrorJSON(c, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("error code 2")
		ErrorJSON(c, errors.New("error calling auth service"+err.Error()))
		return
	}

	//Create variable to read response.Body
	var jsonFromService jsonResponse

	//Decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println("error decode")
		ErrorJSON(c, err)
		return
	}

	if jsonFromService.Error {
		log.Println("error")
		ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	// WriteJSON(c, http.StatusAccepted, payload)
	c.JSON(http.StatusAccepted, payload)
}

func SendMail(c *gin.Context, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mail-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		ErrorJSON(c, errors.New("error calling mail service"))
	}

	//Send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	c.JSON(http.StatusAccepted, payload)
}

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})

}
