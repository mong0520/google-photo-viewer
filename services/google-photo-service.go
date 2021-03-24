package services

import (
    "context"
    "fmt"
    "github.com/mong0520/google-photo-viewer/utils"
    "strings"

    log "github.com/sirupsen/logrus"

    "golang.org/x/oauth2"
    "google.golang.org/api/photoslibrary/v1"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

type GooglePhotoService struct {
    Service *photoslibrary.Service
    AlbumsService *photoslibrary.AlbumsService
    //PeopleService *people.PeopleService
    //RedisService *redis.Client
}
var googlePhotoService *GooglePhotoService

type WrappedGooglePhotoAlbum struct {
    Title string `json:"title"`
    Url   string `json:"url"`
}

type GetGetAlbumsOptions struct {
    AccountIndex string
}

type WrappedGooglePhotoAlbums []WrappedGooglePhotoAlbum

func GetGooglePhotoService(user string, conf *oauth2.Config) (*GooglePhotoService, error){
    if googlePhotoService != nil{
        return googlePhotoService, nil
    }

    googlePhotoService = &GooglePhotoService{}

	//forceToken := false
	//if forceToken{
	//    return nil, errors.New("force renew token")
    //}
	service := &photoslibrary.Service{}

	existToken, err := utils.RetrieveToken(user)
	if err != nil{
	    return nil, err
    }
    client := conf.Client(context.Background(), existToken)

	service, err = photoslibrary.New(client)
	if err != nil{
        return googlePhotoService, err
    }


    googlePhotoService.Service = service
    googlePhotoService.AlbumsService = photoslibrary.NewAlbumsService(service)
    //googlePhotoService.PeopleService = people.NewPeopleService(peopleService)
    return googlePhotoService, nil
}

func (g *GooglePhotoService)GetAlbums(options *GetGetAlbumsOptions) ([]WrappedGooglePhotoAlbum, error) {
    var wrappedGooglePhotoAlbums []WrappedGooglePhotoAlbum
    albumsService := g.AlbumsService
	albumList := albumsService.List()
	ret, err := albumList.PageSize(50).Do()
    albumList.Fields()
	albumList.Do()
	if err != nil {
		log.Fatal(err)
		return wrappedGooglePhotoAlbums, err
	}
	// first time
	for _, album := range ret.Albums {
		if options.AccountIndex != ""{
            album.ProductUrl = strings.Replace(album.ProductUrl, "photos.google.com", "photos.google.com/u/"+ options.AccountIndex, -1)
        }
        fmt.Println(album.Title, album.ProductUrl)
        wrappedGooglePhotoAlbums = append(wrappedGooglePhotoAlbums, WrappedGooglePhotoAlbum{
            Title: album.Title,
            Url:   album.ProductUrl,
        })
	}
	for {
		nextPageToken := ret.NextPageToken
		if nextPageToken == "" {
			break
		}
		ret, err = albumList.PageToken(nextPageToken).PageSize(50).Do()
		if err != nil {
			log.Fatal(err)
			return wrappedGooglePhotoAlbums, err
		}
		for _, album := range ret.Albums {
           if options.AccountIndex != ""{
               album.ProductUrl = strings.Replace(album.ProductUrl, "photos.google.com", "photos.google.com/u/"+ options.AccountIndex, -1)
           }
           fmt.Println(album.Title, album.ProductUrl)
           wrappedGooglePhotoAlbums = append(wrappedGooglePhotoAlbums, WrappedGooglePhotoAlbum{
              Title: album.Title,
              Url:   album.ProductUrl,
          })
		}
	}
	return wrappedGooglePhotoAlbums, nil
}
