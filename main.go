package main

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mong0520/google-photo-viewer/handlers"
	"github.com/mong0520/google-photo-viewer/models"
	"github.com/mong0520/google-photo-viewer/services"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func initializeStorageService(r *gin.Engine) {
	godotenv.Load()
	gob.Register(&models.UserInfo{})
	gob.Register(&oauth2.Config{})
	gob.Register(&oauth2.Token{})

	var err error
	if err = services.InitRedisService(os.Getenv("RedisHostPort")); err != nil {
		log.Fatal(err)
	}

	if err = services.InitMongoService(os.Getenv("MongoDBUri")); err != nil {
		log.Fatal(err)
	}

	if err = services.InitSessionService(r, os.Getenv("RedisHostPort")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	initializeStorageService(r)

	r.LoadHTMLGlob("view/*")
	r.GET("/u/:idx/albums", handlers.AlbumHandler)
	// r.GET("/u/:idx/listMedia", handlers.ListMediaItemsHandler)
	r.GET("/u/:idx", handlers.MainHandler)
	// r.GET("/login/u/:idx", handlers.LoginHandler)
	r.GET("/callback", handlers.CallbackHandler)
	r.GET("/check", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/me", handlers.MeHandler)
	r.GET("/auth", handlers.LoginHandler)
	r.GET("/albums", handlers.AlbumHandler)
	r.GET("/albums/save", handlers.SaveAlbumsHandler)
	// r.GET("/media/init", handlers.InitMediaItemsHandler)

	r.Run(":80")
}
