package services

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/joho/godotenv"
    log "github.com/sirupsen/logrus"

    oauth2ns "github.com/nmrshll/oauth2-noserver"
    "github.com/zalando/go-keyring"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/people/v1"
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

func GetGooglePhotoService(user string) (*GooglePhotoService, error){
    if googlePhotoService != nil{
        return googlePhotoService, nil
    }

    googlePhotoService = &GooglePhotoService{}

    godotenv.Load()
	//user, err := user.Current()
	//if err != nil {
	//	panic(err)
	//}
	// ask the user to authenticate on google in the browser
	conf := &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Scopes:       []string{
		    photoslibrary.PhotoslibraryScope,
		    people.UserEmailsReadScope, // required
		    people.UserinfoProfileScope, // required
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	client := &oauth2ns.AuthorizedClient{}
	// Try to use existing token
	existToken, err := retrieveToken(user)
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
		storeToken(user, client.Token)
	} else {
		// Use existing one
		client = &oauth2ns.AuthorizedClient{
			Client: conf.Client(context.Background(), existToken),
			Token:  existToken,
		}
	}
    //peopleService := &people.Service{}
    //peopleService, err = people.New(client.Client)
    //fmt.Println(err)
    //ret, err := peopleService.People.Get("people/me").PersonFields("emailAddresses").Do()
    //fmt.Println(err)
    //fmt.Println(ret.EmailAddresses)

	service, err = photoslibrary.New(client.Client)
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
