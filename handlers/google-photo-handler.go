package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/services"
    "net/http"
)


func Handler(c *gin.Context){
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
}
