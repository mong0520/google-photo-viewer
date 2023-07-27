package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/google-photo-viewer/services"
	"google.golang.org/api/photoslibrary/v1"
	"net/http"
	"time"
)

func SaveMediaItemsHandler(c *gin.Context) {
	// var mediaItems []*photoslibrary.MediaItem

	svc, err := services.GetGooglePhotoService(c)
	if err != nil {
		c.Error(err)
	}

	go svc.UpsertPhotosToDB()

	c.HTML(http.StatusOK, "media_items.html", nil)
	// mediaItems, err = svc.UpsertPhotosToDB()
	// if err != nil {
	// 	c.Error(err)
	// }
	//
	// for _, item := range mediaItems {
	// 	fmt.Println(item)
	// }
}

func SaveAlbumsHandler(c *gin.Context) {
	svc, err := services.GetGooglePhotoService(c)
	if err != nil {
		c.Error(err)
	}

	// force expire redis data
	redisSvc := services.GetRedisService()
	userInfo := services.GetSessionService().GetUserInfo(c)
	redisSvc.Del(c, userInfo.ID)

	// options := &services.GetGetAlbumsOptions{}
	albums, err := svc.GetAlbums()
	if err != nil {
		c.Error(err)
	}

	svc.UpsertAlbumsToDB(albums)

	albumsBytes, err := json.Marshal(albums)
	if err != nil {
		c.Error(err)
	}
	_, err = redisSvc.SetEX(context.Background(), userInfo.ID, albumsBytes, 86400*time.Second).Result()
	if err != nil {
		c.Error(err)
	}
}

func GetAlbumsHandler(c *gin.Context) {
	redisSvc := services.GetRedisService()
	var albums []photoslibrary.Album

	// Get result from redis
	userInfo := services.GetSessionService().GetUserInfo(c)
	result, err := redisSvc.Get(context.Background(), userInfo.ID).Bytes()
	// svc.ListPhotos()
	if err != nil || len(result) == 0 {
		// not hit cache
		// options := &services.GetGetAlbumsOptions{}
		svc, err := services.GetGooglePhotoService(c)
		if err != nil {
			c.Error(err)
		}

		albums, err = svc.GetAlbums()
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
