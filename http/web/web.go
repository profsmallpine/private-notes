package web

import (
	"context"
	"embed"
	"fmt"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/xy-planning-network/trails/http/ctx"
	m "github.com/xy-planning-network/trails/http/middleware"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/postgres"
	"gorm.io/gorm"
)

//go:embed tmpl/**/*.tmpl
var files embed.FS

const (
	tmplDir      = "tmpl/layout"
	authedTmpl   = tmplDir + "/authenticated_base.tmpl"
	unauthedTmpl = tmplDir + "/unauthenticated_base.tmpl"
)

var (
	errorMessage = fmt.Sprintf(session.ContactUsErr, "profsmallpine@gmail.com")
)

type Controller struct {
	Database postgres.DatabaseService
	Keyring  ctx.KeyRingable

	*gorm.DB
	domain.Procedures
	domain.Services
	*resp.Responder
}

func (c *Controller) GetByID(id uint) (m.User, error) {
	user := &domain.User{}
	if err := c.DB.Preload("Groups").First(user, id).Error; err != nil {
		return nil, err
	}

	return m.User(user), nil
}

// currentUser helps by retrieving the User stored from the provided context..
func (c *Controller) currentUser(ctx context.Context) (*domain.User, error) {
	u, ok := ctx.Value(c.Keyring.CurrentUserKey()).(*domain.User)
	if !ok {
		return nil, domain.ErrNoUser
	}
	return u, nil
}

// session helps by retrieving the session.TrailsSessionable from the provided context.
func (c *Controller) session(ctx context.Context) (session.TrailsSessionable, error) {
	s, ok := ctx.Value(c.Keyring.SessionKey()).(session.TrailsSessionable)
	if !ok {
		return nil, domain.ErrNoSession
	}
	return s, nil
}

// key is a string representing a key used in a context.Context.
type key string

const (
	currentUserCtxKey key = "currentUserCtxKey"
	sessionCtxKey     key = "sessionCtxKey"
)

// Key returns key so it can be used as a key a map[string].
func (k key) Key() string { return string(k) }

// String formats the stringified key with additional contextual information
func (k key) String() string {
	return "http context key: " + string(k)
}
