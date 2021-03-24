package handlers

import "C"
import (
    "encoding/json"
    "fmt"
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/mong0520/google-photo-viewer/utils"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/people/v1"
    "google.golang.org/api/photoslibrary/v1"
    "io/ioutil"
    "net/http"
    "os"
)

type UserInfo struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    GivenName  string `json:"given_name"`
    FamilyName string `json:"family_name"`
    Picture    string `json:"picture"`
    Locale     string `json:"locale"`
}

var (
    // TODO: randomize it
    oauthStateString = "pseudo-random"
    conf *oauth2.Config
    userInfo *UserInfo
)


func LoginHandler(c *gin.Context){
    // get Google App clients and secrets
    godotenv.Load()
    // ask the user to authenticate on google in the browser
    // ref: https://itnext.io/getting-started-with-oauth2-in-go-1c692420e03
    conf = &oauth2.Config{
        ClientID:     os.Getenv("ClientID"),
        ClientSecret: os.Getenv("ClientSecret"),
        RedirectURL:  fmt.Sprintf("%s/callback", os.Getenv("HostUrl")),
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

    url := conf.AuthCodeURL(oauthStateString)
    c.Redirect(http.StatusTemporaryRedirect, url)
}


func getAccessToken(state string, code string) (*UserInfo, *oauth2.Token, error) {
    if state != oauthStateString {
        return nil, nil, fmt.Errorf("invalid oauth state")
    }
    token, err := conf.Exchange(oauth2.NoContext, code)
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

    userInfo = &UserInfo{}
    err = json.Unmarshal(contents, userInfo)
    if err != nil{
        return nil, nil, err
    }

    return userInfo, token, nil
}

func CallbackHandler(c *gin.Context){
    session := sessions.Default(c)
    userInfo, token, err := getAccessToken(c.Query("state"), c.Query("code"))
    if err != nil {
        c.Error(err)
    }

    session.Set("user-info", userInfo)
    session.Set("conf", conf)
    session.Save()
    err = utils.StoreToken(userInfo.ID, token)
    if err != nil {
        c.Error(err)
    }
    c.Redirect(http.StatusTemporaryRedirect, "/photos/1")
}
