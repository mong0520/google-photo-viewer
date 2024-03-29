package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mong0520/google-photo-viewer/handlers"
	"github.com/mong0520/google-photo-viewer/models"
	"github.com/mong0520/google-photo-viewer/services"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"time"
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

func validToken(c *gin.Context) {
	token := services.GetSessionService().GetOAuth2Token(c)
	if token != nil && token.Expiry.Before(time.Now()) {
		fmt.Println("token is expired")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
	}
	c.Next()
}

func main() {
	r := gin.Default()
	initializeStorageService(r)

	r.LoadHTMLGlob("view/*")
	r.GET("/check", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/", validToken, handlers.PortalHandler)
	r.GET("/auth", handlers.AuthHandler)
	r.GET("/callback", handlers.CallbackHandler)
	r.GET("/albums", validToken, handlers.GetAlbumsHandler)
	r.GET("/albums/save", validToken, handlers.SaveAlbumsHandler)
	r.GET("/media_items/save", validToken, handlers.SaveMediaItemsHandler)

	r.Run(":80")
}
