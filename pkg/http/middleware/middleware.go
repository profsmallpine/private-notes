package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Manager struct {
	store sessions.Store
}

func NewManager(authKey, encryptionKey []byte) *Manager {
	store := sessions.NewCookieStore(authKey, encryptionKey)
	// store.Options.Secure = env != domain.Development
	store.Options.HttpOnly = true
	store.MaxAge(86400) // 1 day

	return &Manager{store}
}

func Chain(handler http.Handler, middlewares ...mux.MiddlewareFunc) http.Handler {
	// Loop in reverse to preserve middleware order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
