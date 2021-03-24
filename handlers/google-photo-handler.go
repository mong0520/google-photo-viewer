package handlers

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/services"
    "net/http"
    "time"
)


func Handler(c *gin.Context){
    accountIdxInt := c.Param("idx")
    user := c.GetHeader("x-authorized-user")
    //temporary hard-code
    user = "mong"
    redisSvc := services.GetRedisService()
    if redisSvc == nil{
        c.Error(errors.New("unable to get redis client"))
    }

    svc, err := services.GetGooglePhotoService(user)
    if err != nil {
        c.Error(err)
    }

    options := &services.GetGetAlbumsOptions{AccountIndex: accountIdxInt}

    // Get result from redis
    var albums []services.WrappedGooglePhotoAlbum
    result, err := redisSvc.Get(context.Background(), generateRedisKey(user, accountIdxInt)).Bytes()
    if err != nil || len(result) == 0 {
        albums, err = svc.GetAlbums(options)
        if err != nil {
            c.Error(err)
        }

        albumsBytes, err :=  json.Marshal(albums)
        if err != nil {
            c.Error(err)
        }
        _, err = redisSvc.SetEX(context.Background(), generateRedisKey(user, accountIdxInt), albumsBytes, 86400 * time.Second).Result()
        if err != nil {
            c.Error(err)
        }
    }else{
        err := json.Unmarshal(result, &albums)
        if err != nil {
            c.Error(err)
        }
    }

    c.HTML(http.StatusOK, "albums.html", gin.H{
        "albums": albums,
    })

}

func generateRedisKey(user string, idx string) string{
    return fmt.Sprintf("user:%s:idx:%s", user, idx)
}
