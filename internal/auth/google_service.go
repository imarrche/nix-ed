package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/imarrche/nix-ed/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleService is a service for Google oauth.
type GoogleService struct {
	c *oauth2.Config
}

// NewGoogleService creates and returns a new GoogleService instance.
func NewGoogleService() *GoogleService {
	return &GoogleService{
		c: &oauth2.Config{
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			ClientID:     config.Get().ClientID,
			ClientSecret: config.Get().ClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

// AuthCodeURL returns authentication code URL.
func (s *GoogleService) AuthCodeURL(state string) string {
	return s.c.AuthCodeURL(state)
}

// GetAccessToken returns access token.
func (s *GoogleService) GetAccessToken(code string) (string, error) {
	token, err := s.c.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	return token.AccessToken, nil
}

// GetUserInfo returns the information about user (ID, email and etc).
func (s *GoogleService) GetUserInfo(token string) ([]byte, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
