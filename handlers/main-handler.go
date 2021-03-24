package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)


func MainHandler(c *gin.Context){
    accountIdx := c.Param("idx")
    c.HTML(http.StatusOK, "login.html", gin.H{
        "accountIdx": accountIdx,
    })
}
