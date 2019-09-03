package spotifyapi

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Scopes let the developer specify what type of data they want
// The scopes that the user passes in during authentication determines
// what data the user is asked to grant
const (
	ScopeUserTopRead = "user-top-read"
)

// Authenticator provides functions for implementing the OAuth2 flow.
type Authenticator struct {
	config  *oauth2.Config
	context context.Context
}

//NewAuthenticator creates an authenticator which is used to implement the OAuth2 flow.
func NewAuthenticator(redirectURL string, clientID string, clientSecret string, scopes ...string) Authenticator {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     spotify.Endpoint,
	}

	ctx := context.TODO()

	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

//AuthCodeURL returns a URL to the the Spotify Accounts Service's OAuth2 endpoint.
// The state is a token used to protect the user from CSRF attacks.
//For more info, refer to http://tools.ietf.org/html/rfc6749#section-10.12.
func (a Authenticator) AuthCodeURL(state string, showDialog bool) string {
	if showDialog {
		return a.config.AuthCodeURL(state, oauth2.SetAuthURLParam("show_dialog", "true"))
	}
	return a.config.AuthCodeURL(state)
}

// Token attempts to pull an authorization code from an HTTP request
// and exchange it for an access token
func (a Authenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	parameters := r.URL.Query()
	if err := parameters.Get("error"); err != "" {
		return nil, fmt.Errorf("Spotify authorization failed - " + err)
	}

	code := parameters.Get("code")
	if code == "" {
		return nil, fmt.Errorf("Spotify didn't return an access code")
	}

	returnedState := parameters.Get("state")
	if state != returnedState {
		return nil, fmt.Errorf("Spotify redirect state parameter doesn't match")
	}

	return a.config.Exchange(a.context, code)
}

//NewClient creates an http.client that will be use the specified acccess
//token for its API requests
func (a Authenticator) NewClient(token *oauth2.Token) Client {
	httpClient := a.config.Client(a.context, token)
	newClient := Client{
		client:    httpClient,
		baseURL:   spotifyBaseURL,
		AutoRetry: true,
	}

	return newClient
}
