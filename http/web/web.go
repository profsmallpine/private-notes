package web

import (
	"context"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/ranger"
	"gorm.io/gorm"
)

type Controller struct {
	*gorm.DB
	*ranger.Ranger
	domain.Procedures
	domain.Services
}

// currentUser helps by retrieving the User stored from the provided context..
func (h *Controller) currentUser(ctx context.Context) (*domain.User, error) {
	u, ok := ctx.Value(h.EmitKeyring().CurrentUserKey()).(*domain.User)
	if !ok {
		return nil, domain.ErrNoUser
	}
	return u, nil
}

// session helps by retrieving the session.TrailsSessionable from the provided context.
func (h *Controller) session(ctx context.Context) (session.TrailsSessionable, error) {
	s, ok := ctx.Value(h.EmitKeyring().SessionKey()).(session.TrailsSessionable)
	if !ok {
		return nil, domain.ErrNoSession
	}
	return s, nil
}
