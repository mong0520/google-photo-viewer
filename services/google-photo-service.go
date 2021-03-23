package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	oauth2ns "github.com/nmrshll/oauth2-noserver"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
    keyring "github.com/zalando/go-keyring"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

type GooglePhotoService struct {
    Service *photoslibrary.Service
    AlbumsService *photoslibrary.AlbumsService
}
var googlePhotoService *GooglePhotoService

type WrappedGooglePhotoAlbum struct {
    Title string
    Url   string
}

type WrappedGooglePhotoAlbums []WrappedGooglePhotoAlbum

func GetGooglePhotoService() (*GooglePhotoService, error){
    if googlePhotoService != nil{
        return googlePhotoService, nil
    }

    googlePhotoService = &GooglePhotoService{}

    godotenv.Load()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// ask the user to authenticate on google in the browser
	conf := &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	client := &oauth2ns.AuthorizedClient{}
	// Try to use existing token
	existToken, err := retrieveToken(user.Name)
	forceToken := false
	service := &photoslibrary.Service{}

	if err != nil || forceToken == true {
		// Token not found
		log.Debug(err)

		// Request a new access token
		client, err = oauth2ns.AuthenticateUser(conf)
		if err != nil {
			log.Debug(err)
		}

		// Store it
		storeToken(user.Name, client.Token)
	} else {
		// Use existing one
		client = &oauth2ns.AuthorizedClient{
			Client: conf.Client(context.Background(), existToken),
			Token:  existToken,
		}
	}

	service, err = photoslibrary.New(client.Client)
	if err != nil{
        return googlePhotoService, err
    }

    googlePhotoService.Service = service
    googlePhotoService.AlbumsService = photoslibrary.NewAlbumsService(service)
    return googlePhotoService, nil
}

func (g *GooglePhotoService)GetAlbums() ([]WrappedGooglePhotoAlbum, error) {
    var wrappedGooglePhotoAlbums []WrappedGooglePhotoAlbum
    albumsService := g.AlbumsService
	albumList := albumsService.List()
	ret, err := albumList.PageSize(50).Do()
	albumList.Do()
	if err != nil {
		log.Fatal(err)
		return wrappedGooglePhotoAlbums, err
	}
	for _, album := range ret.Albums {
		fmt.Println(album.Title, album.ProductUrl)
        wrappedGooglePhotoAlbums = append(wrappedGooglePhotoAlbums, WrappedGooglePhotoAlbum{
            Title: album.Title,
            Url:   album.ProductUrl,
        })
	}
	//for {
	//	nextPageToken := ret.NextPageToken
	//	if nextPageToken == "" {
	//		break
	//	}
	//	ret, err = albumList.PageToken(nextPageToken).PageSize(50).Do()
	//	if err != nil {
	//		log.Fatal(err)
	//		return wrappedGooglePhotoAlbums, err
	//	}
	//	for _, album := range ret.Albums {
	//		fmt.Println(album.Title, album.ProductUrl)
    //        wrappedGooglePhotoAlbums = append(wrappedGooglePhotoAlbums, WrappedGooglePhotoAlbum{
    //            Title: album.Title,
    //            Url:   album.ProductUrl,
    //        })
	//	}
	//}

	return wrappedGooglePhotoAlbums, nil
}

func storeToken(googleUserEmail string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, googleUserEmail, string(tokenJSONBytes))
	if err != nil {
		log.Debugf("failed storing token into keyring: %v", err)
		return err
	}

	return nil
}

func retrieveToken(googleUserEmail string) (*oauth2.Token, error) {
	tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, err
		}

		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(tokenJSONString), &token)
	if err != nil {
		log.Debugf("failed unmarshaling token: %v", err)
		return nil, err
	}

	// validate token
	if !token.Valid() {
		return nil, errors.New("invalid token")
	}

	return &token, nil
}
