package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type RPCPayload struct {
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
		Authenticate(c, requestPayload.Auth)
	case "log":
		//LogItem(c, requestPayload.Log)
		//LogEventViaRabbit(c, requestPayload.Log)
		log.Println("enter here to log")
		LogItemViaRPC(c, requestPayload.Log)
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
		ErrorJSON(c, errors.New("error log service"))
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
		log.Println("error request handler", err)
		ErrorJSON(c, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("error response handler", err)
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the correct status
	log.Println("status code :", response.StatusCode)
	if response.StatusCode != http.StatusAccepted {
		log.Println("error response status code")
		ErrorJSON(c, errors.New("invalid credentials"))
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
	log.Println("enter here")
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mail-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error request handler", err)
		ErrorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("error response handler", err)
		ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	//Make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		log.Println("error status code")
		ErrorJSON(c, errors.New("error calling mail service"))
	}

	//Send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	c.JSON(http.StatusAccepted, payload)
}

func LogEventViaRabbit(c *gin.Context, l LogPayload) {
	err := PushToQueue(c, l.Name, l.Data)
	if err != nil {
		log.Println("error push to queue")
		ErrorJSON(c, err)
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	c.JSON(http.StatusAccepted, payload)
}

func PushToQueue(c *gin.Context, name, msg string) error {
	rabbitConn, err := Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	emitter, err := event.NewEventEmitter(rabbitConn)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(c, string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

func LogItemViaRPC(c *gin.Context, l LogPayload) {
	log.Println("entered log item via rpc")
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Println("error in client")
		ErrorJSON(c, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	log.Println("client: ", client)

	var result string
	log.Println("rpc Payload: ", rpcPayload)
	//log.Println("error is: ", client.Call("RPCServer.LogInfo", rpcPayload, &result))
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		log.Println("error in result")
		ErrorJSON(c, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	c.JSON(http.StatusAccepted, payload)
}

func LogViaGRPC(c *gin.Context) {
	log.Println("enter broker log via gRPC")
	var requestPayload RequestPayload

	err := ReadJSON(c, &requestPayload)
	if err != nil {
		ErrorJSON(c, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("error in dial gRPC")
		ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	client := logs.NewLogServiceClient(conn)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	_, err = client.WriteLog(c, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		log.Println("error in write log gRPC")
		ErrorJSON(c, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged via gRPC"

	c.JSON(http.StatusAccepted, payload)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})

}
