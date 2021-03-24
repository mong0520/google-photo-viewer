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
        c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/%s", accountIdxInt))
    }
    userInfo := userInfoVal.(*UserInfo)
    confVal := session.Get("conf")
    if confVal == nil {
        c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/%s", accountIdxInt))
    }
    conf := confVal.(*oauth2.Config)

    tokenVal := session.Get("token")
    if tokenVal == nil {
        c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/%s", accountIdxInt))
    }
    token := tokenVal.(*oauth2.Token)
    redisSvc := services.GetRedisService()
    if redisSvc == nil{
        c.Error(errors.New("unable to get redis client"))
    }

    // Get result from redis
    var albums []services.WrappedGooglePhotoAlbum
    result, err := redisSvc.Get(context.Background(), userInfo.ID).Bytes()
    if err != nil || len(result) == 0 {
        // not hit cache
        options := &services.GetGetAlbumsOptions{AccountIndex: accountIdxInt}
        svc, err := services.GetGooglePhotoService(conf, token)
        if err != nil {
            c.Error(err)
        }
        albums, err = svc.GetAlbums(options)
        if err != nil {
            c.Error(err)
        }

        albumsBytes, err :=  json.Marshal(albums)
        if err != nil {
            c.Error(err)
        }
        _, err = redisSvc.SetEX(context.Background(), userInfo.ID, albumsBytes, 86400 * time.Second).Result()
        if err != nil {
            c.Error(err)
        }
    }else{
        // hit cache
        err := json.Unmarshal(result, &albums)
        if err != nil {
            c.Error(err)
        }
    }

    c.HTML(http.StatusOK, "albums.html", gin.H{
        "albums": albums,
    })

}
