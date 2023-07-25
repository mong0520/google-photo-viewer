package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/google-photo-viewer/services"
	"github.com/mong0520/google-photo-viewer/utils"
	"google.golang.org/api/photoslibrary/v1"
	"net/http"
	"time"
)

func InitMediaItemsHandler(c *gin.Context) {
	conf, err := utils.RetrieveOAuthConf(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"errorMessage": err,
		})
	}
	token, err := utils.RetrieveOAuthToken(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"errorMessage": err,
		})
	}

	var mediaItems []*photoslibrary.MediaItem

	svc, err := services.GetGooglePhotoService(conf, token)
	if err != nil {
		c.Error(err)
	}

	mediaItems, err = svc.UpsertPhotosToDB()
	if err != nil {
		c.Error(err)
	}

	for _, item := range mediaItems {
		fmt.Println(item)
	}
}

func AlbumHandler(c *gin.Context) {
	session := sessions.Default(c)
	accountIdxInt := c.Param("idx")
	userInfoVal := session.Get("user-info")
	if userInfoVal == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	userInfo := userInfoVal.(*UserInfo)
	confVal := session.Get("conf")
	if confVal == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	conf, err := utils.RetrieveOAuthConf(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"errorMessage": err,
		})
	}
	token, err := utils.RetrieveOAuthToken(c)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"errorMessage": err,
		})
	}

	redisSvc := services.GetRedisService()
	if redisSvc == nil {
		c.Error(errors.New("unable to get redis client"))
		return
	}

	// Get result from redis
	var albums []services.WrappedGooglePhotoAlbum
	result, err := redisSvc.Get(context.Background(), userInfo.ID).Bytes()
	// svc.ListPhotos()
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

		albumsBytes, err := json.Marshal(albums)
		if err != nil {
			c.Error(err)
		}
		_, err = redisSvc.SetEX(context.Background(), userInfo.ID, albumsBytes, 86400*time.Second).Result()
		if err != nil {
			c.Error(err)
		}
	} else {
		// hit cache
		err := json.Unmarshal(result, &albums)
		if err != nil {
			c.Error(err)
		}
	}

	c.HTML(http.StatusOK, "albums.html", gin.H{
		"albums":   albums,
		"userInfo": userInfo,
	})

}
