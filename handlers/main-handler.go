package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)


func MainHandler(c *gin.Context){
    c.HTML(http.StatusOK, "login.html", nil)
}
