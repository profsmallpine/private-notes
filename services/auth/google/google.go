package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/profsmallpine/private-notes/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type Service struct {
	Config *oauth2.Config
}

func NewService(callbackURL string) domain.AuthProvider {
	return Service{
		Config: &oauth2.Config{
			RedirectURL:  os.Getenv("BASE_URL") + callbackURL,
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			Scopes:       []string{goauth2.UserinfoEmailScope, goauth2.UserinfoProfileScope},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s Service) ExchangeCode(r *http.Request) (*domain.AuthData, error) {
	// Read oauthState from Cookie and check against state param.
	oauthState, _ := r.Cookie("oauthstate")
	if oauthState == nil || r.FormValue("state") != oauthState.Value {
		return nil, domain.ErrUnauthorized
	}

	// Use code to get token.
	token, err := s.Config.Exchange(context.Background(), r.FormValue("code"), oauth2.AccessTypeOffline)
	if err != nil {
		return nil, err
	}

	// Create oauth2 service
	ctx := context.Background()
	service, err := goauth2.NewService(ctx, option.WithTokenSource(s.Config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	// Fetch user data with oauth2 service
	user, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &domain.AuthData{
		Email:      user.Email,
		FirstName:  user.GivenName,
		LastName:   user.FamilyName,
		PictureURL: user.Picture,
	}, nil
}

func (s Service) FetchAuthURL(w http.ResponseWriter) (string, error) {
	// Create oauthState cookie
	expiration := time.Now().Add(time.Hour)
	b := make([]byte, 16)
	_, err := rand.Read(make([]byte, 16))
	if err != nil {
		return "", err
	}

	oauthState := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: oauthState, Expires: expiration}
	http.SetCookie(w, &cookie)

	// AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
	// You must always provide a non-empty string and validate that it matches the the
	// state query parameter on your redirect callback.
	return s.Config.AuthCodeURL(oauthState, oauth2.ApprovalForce), nil
}
