package github

import (
	"github.com/gu-io/midash/pkg/internals/auth"
	"golang.org/x/oauth2/github"
)

// New returns a new instance of auth.Auth for use with the google OAuth2 API.
func New(cred auth.Credential, redirectURL string) *auth.Auth {
	return auth.New(cred, github.Endpoint, redirectURL)
}
