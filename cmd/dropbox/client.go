package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func createClient(appKey, appSecret, domain string) (*http.Client, *oauth2.Token, error) {
	conf := &oauth2.Config{
		ClientID:     appKey,
		ClientSecret: appSecret,
		Endpoint:     dropbox.OAuthEndpoint(domain),
		/*
			oauth2.Endpoint{
				AuthURL:  "https://www.dropbox.com/oauth2/authorize",
				TokenURL: "https://api.dropboxapi.com/oauth2/token",
			},
			Scopes:       nil,
		*/
	}

	// Read token from file
	var tok *oauth2.Token
	executable, err := os.Executable()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get executable")
	}
	tokenFile := executable + ".token"
	tokdat, err := ioutil.ReadFile(tokenFile)
	if err == nil {
		tok = &oauth2.Token{}
		if err := json.Unmarshal(tokdat, tok); err != nil {
			return nil, nil, errors.Wrap(err, "unable to deserialize token from JSON")
		}
	} else if os.IsNotExist(err) {
		url := conf.AuthCodeURL(uuid.New().String())
		fmt.Printf("Visit this URL: %v\n", url)
		fmt.Print("Enter the access code here: ")
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			return nil, nil, errors.Wrap(err, "unable to scan code")
		}
		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable for token exchange")
		}
		fmt.Printf("Token: %v\n", tok)
		b, err := json.Marshal(tok)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to serialize token to JSON")
		}
		if err := ioutil.WriteFile(tokenFile, b, 0600); err != nil {
			return nil, nil, errors.Wrap(err, "unable to write token file")
		}
	} else {
		return nil, nil, errors.Wrap(err, "unable to open token file")
	}

	client := conf.Client(oauth2.NoContext, tok)

	return client, tok, nil
}
