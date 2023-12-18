package main

import (
	"log-service/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})

}

func WriteLog(c *gin.Context) {
	//Add model to writelog
	mdl := data.New(client).LogEntry

	//Read Json into var
	var requestPayload JSONPayload
	_ = ReadJSON(c, &requestPayload)

	//Insert Data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := mdl.Insert(event)
	if err != nil {
		ErrorJSON(c, err)
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	//WriteJSON(c, http.StatusAccepted, resp)
	c.JSON(http.StatusAccepted, resp)
}
