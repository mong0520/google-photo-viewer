package handlers

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/mong0520/google-photo-viewer/services"
    "golang.org/x/oauth2"
    "net/http"
    "time"
)


func GooglePhotoHandler(c *gin.Context){
    session := sessions.Default(c)
    accountIdxInt := c.Param("idx")
    userInfoVal := session.Get("user-info")
    if userInfoVal == nil {
        c.Redirect(http.StatusTemporaryRedirect, "/")
    }
    userInfo := userInfoVal.(*UserInfo)
    confVal := session.Get("conf")
    if confVal == nil {
        c.Redirect(http.StatusTemporaryRedirect, "/")
    }
    conf := confVal.(*oauth2.Config)
    redisSvc := services.GetRedisService()
    if redisSvc == nil{
        c.Error(errors.New("unable to get redis client"))
    }

    svc, err := services.GetGooglePhotoService(userInfo.ID, conf)
    if err != nil {
        c.Error(err)
    }

    options := &services.GetGetAlbumsOptions{AccountIndex: accountIdxInt}

    // Get result from redis
    var albums []services.WrappedGooglePhotoAlbum
    result, err := redisSvc.Get(context.Background(), generateRedisKey(userInfo.ID, accountIdxInt)).Bytes()
    if err != nil || len(result) == 0 {
        albums, err = svc.GetAlbums(options)
        if err != nil {
            c.Error(err)
        }

        albumsBytes, err :=  json.Marshal(albums)
        if err != nil {
            c.Error(err)
        }
        _, err = redisSvc.SetEX(context.Background(), generateRedisKey(userInfo.ID, accountIdxInt), albumsBytes, 86400 * time.Second).Result()
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
