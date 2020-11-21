package middleware

import (
	"fmt"
	"net/http"
)

const (
	currentUserKey        = "userID"
	sessionKey            = "session"
	currentUserCtxKey key = "currentUserCtxKey"
)

type key string

func (k key) String() string {
	return "pkg/http/middleware context key: " + string(k)
}

func (m *Manager) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := m.store.Get(r, sessionKey)
	if err != nil {
		// TODO: log error
		fmt.Println(err)
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		// TODO: log error
		fmt.Println(err)
	}
}

func (m *Manager) RegisterSession(w http.ResponseWriter, r *http.Request, userID uint) {
	session, err := m.store.Get(r, sessionKey)
	if err != nil {
		// TODO: log error
		fmt.Println(err)
	}

	session.Values[currentUserKey] = userID
	if err := session.Save(r, w); err != nil {
		// TODO: log and redirect back to login if error
		fmt.Println(err)
	}
}
