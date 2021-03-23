package main

import (
    "github.com/mong0520/google-photo-viewer/services"
    "log"
    "net/http"

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

    //svc, err := services.InitGooglePhotoService()
    //if err != nil{
    //    panic(err)
    //}
    books := make([]Book, 0)
    books = append(books, Book{
        Title:  "Title 1",
        Author: "Author 1",
    })
    books = append(books, Book{
        Title:  "Title 2",
        Author: "Author 2",
    })

    r.GET("/demo", func(c *gin.Context) {
        c.HTML(http.StatusOK, "demo.html", gin.H{
            "books": books,
        })
    })
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
