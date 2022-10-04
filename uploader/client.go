package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/omarahm3/automato/uploader/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetClient(scope string, config *config.Config) (*http.Client, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(config.SecretsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets file %v", err)
	}

	c, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secrets file to config %v", err)
	}

	token, err := tokenFromFile(config.TokenFile)
	if err != nil {
		token, err = getTokenFromPrompt(c)
		if err != nil {
			return nil, err
		}

		err = saveToken(config.TokenFile, token)
		if err != nil {
			return nil, err
		}
	}

	return c.Client(ctx, token), nil
}

func saveToken(file string, token *oauth2.Token) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

func getTokenFromPrompt(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("unable to read authorization code %v", err)
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}

	return token, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}
