package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(c *gin.Context, data any) error {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, data)
	log.Println("data :", data)
	if err != nil {
		log.Println("masuk error unmarshal :", err.Error())
		return err
	}

	return nil
}

func WriteJSON(c *gin.Context, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Request.Header[key] = value
		}
	}

	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Access-Control-Allow-Headers", "*")

	c.Writer.WriteHeader(status)

	_, err = c.Writer.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(c *gin.Context, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(c, statusCode, payload)
}
