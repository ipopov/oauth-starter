package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	displaySecret = "urn:ietf:wg:oauth:2.0:oob"
	apiScope      = "https://www.googleapis.com/auth/tasks.readonly"
)

func createTokenAndStore(conf *oauth2.Config, f string) error {
	url := conf.AuthCodeURL("state")
	log.Printf("Please visit\n%s\nand enter it here: ", url)
	var authCode string
	fmt.Scanln(&authCode)
	tok, err := conf.Exchange(oauth2.NoContext, authCode)
	s, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, s, 0644)
}

type Secrets struct {
	ClientID     string
	ClientSecret string
}

func unmarshalSecrets(f string) (*Secrets, error) {
	s, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var ret Secrets
	json.Unmarshal(s, &ret)
	return &ret, nil
}

func main() {
	tokenFile := flag.String("token-file", "", "")
	secretsFile := flag.String("secrets-file", "", "")
	flag.Parse()

	secrets, err := unmarshalSecrets(*secretsFile)
	if err != nil {
		log.Panic(err)
	}

	conf := &oauth2.Config{
		ClientID:     secrets.ClientID,
		ClientSecret: secrets.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			apiScope,
		},
		RedirectURL: displaySecret,
	}

	if _, err := os.Stat(*tokenFile); os.IsNotExist(err) {
		createTokenAndStore(conf, *tokenFile)
	}

	s, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		log.Panic(err)
	}
	var seedToken oauth2.Token
	json.Unmarshal(s, &seedToken)
	currentToken, err := conf.TokenSource(oauth2.NoContext, &seedToken).Token()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Access token:\n%s\n", currentToken.AccessToken)
	fmt.Printf("Authorization header:\nAuthorization: Bearer %s\n", currentToken.AccessToken)
}
