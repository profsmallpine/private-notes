package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/pkg/http/routes"
)

func (m *Manager) GetCurrentUserID(r *http.Request) uint {
	return r.Context().Value(currentUserCtxKey).(uint)
}

func (m *Manager) LoadCurrentUserID() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := m.store.Get(r, sessionKey)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			id, ok := session.Values[currentUserKey].(uint)
			if id == 0 || !ok {
				h.ServeHTTP(w, r)
				return
			}

			// Save session to update cookie expiration
			if err := session.Save(r, w); err != nil {
				// TODO: log if save fails
				return
			}

			ctx := context.WithValue(r.Context(), currentUserCtxKey, id)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *Manager) RequireAuthenticatedUser() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := r.Context().Value(currentUserCtxKey).(uint)
			if id == 0 || !ok {
				switch r.Header.Get("Content-type") {
				case "application/json":
					w.WriteHeader(http.StatusUnauthorized)
					return
				default:
					redirect := routes.GetLoginURL

					if r.Method == http.MethodGet && r.URL.Path != routes.GetLoginURL {
						redirect += "?next=" + r.URL.Path
					}

					http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

func (m *Manager) RequireUnauthenticatedUser() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, _ := r.Context().Value(currentUserCtxKey).(uint)
			if id != 0 {
				http.Redirect(w, r, routes.GetGroupsURL, http.StatusTemporaryRedirect)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
