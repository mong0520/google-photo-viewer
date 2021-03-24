package main

import (
    "encoding/gob"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/handlers"
    "golang.org/x/oauth2"
)

// Book ...
type Book struct {
	Title  string
	Author string
}



func main() {
    gob.Register(&handlers.UserInfo{})
    gob.Register(&oauth2.Config{})
	r := gin.Default()
    store := cookie.NewStore([]byte("secret"))
    r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("view/*")

	r.GET("/photos/:idx", handlers.GooglePhotoHandler)
    r.GET("/", handlers.MainHandler)
    r.GET("/login", handlers.LoginHandler)
    r.GET("/callback", handlers.CallbackHandler)
    r.Run(":8080")
}
