package auth

import (
	"net/http"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/auth/google"
)

var getAuthProvider = map[string]func(string) domain.AuthProvider{
	domain.GoogleAuth: google.NewService,
}

type Service struct {
	baseURL string
}

func NewService(baseURL string) *Service {
	return &Service{baseURL: baseURL}
}

func (s *Service) FetchAuthURL(provider, callbackPath string, w http.ResponseWriter) (string, error) {
	authProvider, ok := getAuthProvider[provider]
	if !ok {
		return "", domain.ErrNotValid
	}
	return authProvider(s.baseURL + callbackPath).FetchAuthURL(w)
}

func (s *Service) ExchangeCode(provider, callbackPath string, r *http.Request) (*domain.AuthData, error) {
	authProvider, ok := getAuthProvider[provider]
	if !ok {
		return nil, domain.ErrNotValid
	}
	return authProvider(s.baseURL + callbackPath).ExchangeCode(r)
}
