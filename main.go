package main

import (
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/handlers"
    "log"
)

// Book ...
type Book struct {
	Title  string
	Author string
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("view/*")

	r.GET("/photos/:idx", handlers.Handler)
	log.Fatal(r.Run())
}
