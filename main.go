package main

import (
    "encoding/gob"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/handlers"
    "golang.org/x/oauth2"
    "net/http"
)

// Book ...
type Book struct {
	Title  string
	Author string
}



func main() {
    gob.Register(&handlers.UserInfo{})
    gob.Register(&oauth2.Config{})
    gob.Register(&oauth2.Token{})
	r := gin.Default()
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("view/*")

	r.GET("/u/:idx/albums", handlers.GooglePhotoHandler)
    r.GET("/u/:idx", handlers.MainHandler)
    r.GET("/login/u/:idx", handlers.LoginHandler)
    r.GET("/callback", handlers.CallbackHandler)
    r.GET("/check", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })
    r.Run(":8080")
}
