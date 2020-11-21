package domain

import "net/http"

const GoogleAuth = "google"

type AuthData struct {
	Email      string
	FirstName  string
	LastName   string
	PictureURL string
}

type AuthService interface {
	ExchangeCode(provider, callbackPath string, r *http.Request) (*AuthData, error)
	FetchAuthURL(provider, callbackPath string, w http.ResponseWriter) (string, error)
}

type AuthProvider interface {
	ExchangeCode(r *http.Request) (*AuthData, error)
	FetchAuthURL(w http.ResponseWriter) (string, error)
}
