package main

import (
	"log"
	"net/http"

	"github.com/mong0520/google-photo-viewer/services"

	"github.com/gin-gonic/gin"
)

// Book ...
type Book struct {
	Title  string
	Author string
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("view/*")

	r.GET("/", func(c *gin.Context) {
		svc, err := services.GetGooglePhotoService()
		if err != nil {
			c.Error(err)
		}

		albums, err := svc.GetAlbums()
		if err != nil {
			c.Error(err)
		}

		c.HTML(http.StatusOK, "albums.html", gin.H{
			"albums": albums,
		})
	})
	log.Fatal(r.Run())
}
