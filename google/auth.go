package google

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vcoder4c/jirakits/common"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	credentialFile string = "credential.json"
	tokenFile      string = "token.json"
)

func GetCredentialPath(projectDir string) string {
	return path.Join(projectDir, credentialFile)
}

func getTokenPath(projectDir string) string {
	return path.Join(projectDir, tokenFile)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func getTokenFromFile(projectDir string) (*oauth2.Token, error) {
	f, err := os.Open(getTokenPath(projectDir))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Saves a token to a file path.
func saveToken(projectDir string, token *oauth2.Token) {
	f, err := os.OpenFile(getTokenPath(projectDir), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(token)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(projectDir string, config *oauth2.Config) *http.Client {
	token, err := getTokenFromFile(projectDir)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(projectDir, token)
	}
	return config.Client(context.Background(), token)
}

func IsReady(projectDir string) error {
	credentialPath := GetCredentialPath(projectDir)
	if !common.FileExists(credentialPath) {
		return CredentialIsNotExisted
	}
	tokenPath := getTokenPath(projectDir)
	if !common.FileExists(tokenPath) {
		return TokenIsNotExisted
	}
	_, err := getTokenFromFile(projectDir)
	if err != nil {
		return TokenIsInvalid
	}
	return nil
}
