package handlers

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/google-photo-viewer/models"
	"github.com/mong0520/google-photo-viewer/services"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
	"google.golang.org/api/photoslibrary/v1"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
	conf             *oauth2.Config
)

func MeHandler(c *gin.Context) {
	userInfo := services.GetSessionService().GetUserInfo(c)
	if userInfo == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	if userInfo != nil {
		c.HTML(http.StatusOK, "me.html", gin.H{
			"userInfo": userInfo,
		})
	}
}

func LoginHandler(c *gin.Context) {
	// godotenv.Load()
	// ask the user to authenticate on google in the browser
	// ref: https://itnext.io/getting-started-with-oauth2-in-go-1c692420e03
	conf = &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		RedirectURL:  fmt.Sprintf("%s/callback", os.Getenv("HostUrl")),
		Scopes: []string{
			photoslibrary.PhotoslibraryScope,
			people.UserEmailsReadScope,  // required
			people.UserinfoProfileScope, // required
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}

	url := conf.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackHandler(c *gin.Context) {
	userInfo, token, err := getAccessToken(c.Query("state"), c.Query("code"))
	if err != nil {
		c.Error(err)
	}
	services.GetSessionService().SetSessionEntity(c, "user-info", userInfo)
	services.GetSessionService().SetSessionEntity(c, "conf", conf)
	services.GetSessionService().SetSessionEntity(c, "token", token)

	c.Redirect(http.StatusTemporaryRedirect, "me")
}

func getAccessToken(state string, code string) (*models.UserInfo, *oauth2.Token, error) {
	if state != oauthStateString {
		return nil, nil, fmt.Errorf("invalid oauth state")
	}
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	userInfo := &models.UserInfo{}
	err = json.Unmarshal(contents, userInfo)
	if err != nil {
		return nil, nil, err
	}

	return userInfo, token, nil
}
