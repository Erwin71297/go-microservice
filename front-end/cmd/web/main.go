package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		render(c, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081")
	err := r.Run(":8081")
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(c *gin.Context, t string) {

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var data struct {
		BrokerURL string
	}

	data.BrokerURL = os.Getenv("BROKER_URL")

	if err := tmpl.Execute(c.Writer, data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
