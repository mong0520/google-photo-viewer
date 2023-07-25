package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"google.golang.org/api/photoslibrary/v1"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

var AllPhotos []*photoslibrary.MediaItem

type GooglePhotoService struct {
	Service          *photoslibrary.Service
	AlbumsService    *photoslibrary.AlbumsService
	MediaItemService *photoslibrary.MediaItemsService
	// PeopleService *people.PeopleService
	// RedisService *redis.Client
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

func GetGooglePhotoService(conf *oauth2.Config, token *oauth2.Token) (*GooglePhotoService, error) {
	if googlePhotoService != nil {
		return googlePhotoService, nil
	}

	googlePhotoService = &GooglePhotoService{}
	service := &photoslibrary.Service{}

	// existToken, err := utils.RetrieveToken(user)
	// if err != nil{
	//    return nil, err
	// }

	client := conf.Client(context.Background(), token)
	service, err := photoslibrary.New(client)
	if err != nil {
		return googlePhotoService, err
	}

	googlePhotoService.Service = service
	googlePhotoService.AlbumsService = service.Albums
	googlePhotoService.MediaItemService = service.MediaItems

	// googlePhotoService.PeopleService = people.NewPeopleService(peopleService)
	return googlePhotoService, nil
}

func (g *GooglePhotoService) GetAlbums(options *GetGetAlbumsOptions) ([]WrappedGooglePhotoAlbum, error) {
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
		if options.AccountIndex != "" {
			album.ProductUrl = strings.Replace(album.ProductUrl, "photos.google.com", "photos.google.com/u/"+options.AccountIndex, -1)
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
			if options.AccountIndex != "" {
				album.ProductUrl = strings.Replace(album.ProductUrl, "photos.google.com", "photos.google.com/u/"+options.AccountIndex, -1)
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

func upsertMedia(p *photoslibrary.SearchMediaItemsResponse) error {
	mongodbClient := GetMongoService()
	collection := mongodbClient.Database("google_photos").Collection("mediaItems")

	var pagedItems []photoslibrary.MediaItem
	for _, item := range p.MediaItems {
		pagedItems = append(pagedItems, *item)
	}
	fmt.Println("ready to insert media items to mongodb...")
	var docs []interface{}
	var existingIDs []string
	for _, data := range p.MediaItems {
		docs = append(docs, *data)
		existingIDs = append(existingIDs, data.Id)
	}

	options := options.InsertMany().SetOrdered(false)
	var existingMap map[string]bool
	if len(existingIDs) > 0 {
		// Prepare a filter to check for existing documents
		filter := bson.M{"id": bson.M{"$in": existingIDs}}
		// Find the existing documents in the collection
		existingCursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			log.Fatal("Failed to find existing documents:", err)
		}
		defer existingCursor.Close(context.Background())
		// Prepare a map to store the existing document IDs
		existingMap = make(map[string]bool)
		for existingCursor.Next(context.Background()) {
			var existingData photoslibrary.MediaItem
			if err := existingCursor.Decode(&existingData); err != nil {
				log.Fatal("Failed to decode existing document:", err)
			}
			existingMap[existingData.Id] = true
		}
	}

	// Prepare a list to store the new documents to insert
	var newDocuments []interface{}
	for _, data := range p.MediaItems {
		// Skip the document if it already exists in the collection
		if existingMap[data.Id] {
			fmt.Printf("Document with ID '%s' already exists. Skipping insertion.\n", data.Id)
			continue
		}
		newDocuments = append(newDocuments, data)
	}

	if len(newDocuments) == 0 {
		return nil
	}
	// Insert the JSON documents into MongoDB using InsertMany
	_, err := collection.InsertMany(context.Background(), newDocuments, options)
	if err != nil {
		log.Fatal("Failed to insert documents:", err)
	}
	_, err = collection.InsertMany(context.Background(), docs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Insert %d records to db\n", len(docs))
	}

	return nil
}

func (g *GooglePhotoService) UpsertPhotosToDB() ([]*photoslibrary.MediaItem, error) {
	options := photoslibrary.SearchMediaItemsRequest{PageSize: 100}
	mediaItems := g.Service.MediaItems.Search(&options)
	mediaItems.Do()
	err := mediaItems.Pages(context.Background(), upsertMedia)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
