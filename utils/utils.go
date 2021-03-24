package utils

import (
    "encoding/json"
    "errors"
    "github.com/zalando/go-keyring"
    "golang.org/x/oauth2"
    "log"
)

var (
    serviceName = "google-photo-viewer"
)

func RetrieveToken(googleUserEmail string) (*oauth2.Token, error) {
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
        log.Printf("failed unmarshaling token: %v", err)
        return nil, err
    }

    // validate token
    if !token.Valid() {
        return nil, errors.New("invalid token")
    }

    return &token, nil
}


func StoreToken(googleUserEmail string, token *oauth2.Token) error {
    tokenJSONBytes, err := json.Marshal(token)
    if err != nil {
        return err
    }

    err = keyring.Set("google-photo-viewer", googleUserEmail, string(tokenJSONBytes))
    if err != nil {
        log.Printf("failed storing token into keyring: %v\n", err)
        return err
    }

    return nil
}
