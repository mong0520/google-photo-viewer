package utils

var (
	serviceName = "google-photo-viewer"
)

// func RetrieveOAuthToken(c *gin.Context) (*oauth2.Token, error) {
// 	session := sessions.Default(c)
// 	tokenVal := session.Get("token")
// 	if tokenVal == nil {
// 		return nil, errors.New("unable to find token value of session")
// 	}
// 	token := tokenVal.(*oauth2.Token)
// 	fmt.Println(token.AccessToken)
// 	return token, nil
// }
//
// func RetrieveOAuthConf(c *gin.Context) (*oauth2.Config, error) {
// 	session := sessions.Default(c)
// 	accountIdxInt := c.Param("idx")
// 	userInfoVal := session.Get("user-info")
// 	if userInfoVal == nil {
// 		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/u/%s", accountIdxInt))
// 	}
// 	// userInfo := userInfoVal.(*UserInfo)
// 	confVal := session.Get("conf")
// 	if confVal == nil {
// 		return nil, errors.New("unable to find conf value of session")
// 	}
// 	conf := confVal.(*oauth2.Config)
// 	return conf, nil
// }

// func RetrieveToken(googleUserEmail string) (*oauth2.Token, error) {
//    tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
//    if err != nil {
//        if err == keyring.ErrNotFound {
//            return nil, err
//        }
//
//        return nil, err
//    }
//
//    var token oauth2.Token
//    err = json.Unmarshal([]byte(tokenJSONString), &token)
//    if err != nil {
//        log.Printf("failed unmarshaling token: %v", err)
//        return nil, err
//    }
//
//    // validate token
//    if !token.Valid() {
//        return nil, errors.New("invalid token")
//    }
//
//    return &token, nil
// }

// func StoreToken(googleUserEmail string, token *oauth2.Token) error {
//    tokenJSONBytes, err := json.Marshal(token)
//    if err != nil {
//        return err
//    }
//
//    err = keyring.Set("google-photo-viewer", googleUserEmail, string(tokenJSONBytes))
//    if err != nil {
//        log.Printf("failed storing token into keyring: %v\n", err)
//        return err
//    }
//
//    return nil
// }
