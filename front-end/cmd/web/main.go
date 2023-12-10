package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		render(c, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 80")
	err := r.Run(":80")
	if err != nil {
		log.Panic(err)
	}
}

func render(c *gin.Context, t string) {

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(c.Writer, nil); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}
