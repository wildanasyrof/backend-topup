package oauth

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// type GoogleOauthPkg interface {
// 	AuthCodeURL(state string) string
// 	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
// }

type GoogleOauthPkg struct {
	Config *oauth2.Config
}

func NewGoogleOauthPkg(cfg *config.Config) *GoogleOauthPkg {
	return &GoogleOauthPkg{
		Config: &oauth2.Config{
			ClientID:     cfg.Oauth.ClientID,
			ClientSecret: cfg.Oauth.ClientSecret,
			RedirectURL:  cfg.Oauth.CallbackUrl,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

// AuthCodeURL implements GoogleOauthPkg.
func (g *GoogleOauthPkg) AuthCodeURL(state string) string {
	return g.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange implements GoogleOauthPkg.
func (g *GoogleOauthPkg) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.Config.Exchange(ctx, code)
}
