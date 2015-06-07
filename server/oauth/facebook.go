package oauth

import (
	"fmt"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const facebookAuthCallbackURL = "auth/callback"

var facebookAuthScope = []string{
	"public_profile",
}

type authKey string

func WithFacebook(ctx context.Context) context.Context {
	fbAppID := os.Getenv("FACEBOOK_APPID")
	if fbAppID == "" {
		// TODO:
	}
	fbSecret := os.Getenv("FACEBOOK_SECRET")
	if fbSecret == "" {
		// TODO:
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000/"
	}

	callBackURL := baseURL + facebookAuthCallbackURL

	fmt.Println(callBackURL)

	conf := &oauth2.Config{
		ClientID:     fbAppID,
		ClientSecret: fbSecret,
		RedirectURL:  callBackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/dialog/oauth",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
		Scopes: facebookAuthScope,
	}
	return context.WithValue(ctx, authKey("facebook"), conf)
}

func Facebook(ctx context.Context) *oauth2.Config {
	conf, _ := ctx.Value(authKey("facebook")).(*oauth2.Config)
	return conf
}

func GetFacebookAuthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	c := Facebook(ctx)
	return c.Exchange(oauth2.NoContext, code)
}
