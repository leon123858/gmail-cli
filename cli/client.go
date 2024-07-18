package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	redirectURL = "http://localhost:8080/callback"
)

var mtx sync.Mutex

func getClient(email string) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     viper.Get("id").(string),
		ClientSecret: viper.Get("secret").(string),
		RedirectURL:  redirectURL,
		Scopes:       []string{gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	tokenFile := filepath.Join(GetTokenDir(), fmt.Sprintf("%s.json", email))
	token, err := tokenFromFile(tokenFile)
	if err == nil {
		return config.Client(context.Background(), token), nil
	}

	token, err = getTokenFromWeb(config)
	if err != nil {
		return nil, err
	}
	saveToken(tokenFile, token)

	return config.Client(context.Background(), token), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser: \n%v\n", authURL)

	openBrowser(authURL)

	code := startServerAndWaitForCode()
	if code == "" {
		return nil, fmt.Errorf("didn't get authorization code")
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return token, nil
}

func startServerAndWaitForCode() string {
	var authCode string
	codeChan := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeChan <- code
		_, err := fmt.Fprintf(w, "Authorization successful! You can close this window now.")
		if err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
	})
	if mtx.TryLock() {
		go func() {
			if err := http.ListenAndServe("localhost:8080", nil); err != nil {
				log.Printf("Error starting server: %v", err)
			}
		}()
	}

	authCode = <-codeChan
	return authCode
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Unable to close file: %v", err)
		}
	}(f)
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Unable to close file: %v", err)
		}
	}(f)
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("Unable to encode token: %v", err)
		return
	}
}
