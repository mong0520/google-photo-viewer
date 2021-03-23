package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/services"
    "net/http"
)


func Handler(c *gin.Context){
    accountIdxInt := c.Param("idx")

    svc, err := services.GetGooglePhotoService()
    if err != nil {
        c.Error(err)
    }

    options := &services.GetGetAlbumsOptions{AccountIndex: accountIdxInt}
    albums, err := svc.GetAlbums(options)
    if err != nil {
        c.Error(err)
    }

    c.HTML(http.StatusOK, "albums.html", gin.H{
        "albums": albums,
    })
}
