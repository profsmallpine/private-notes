package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/logger"
)

func (c *Controller) oauthLogin(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	callbackURL := "/auth/" + provider + "/callback"
	rt, err := c.Auth.FetchAuthURL(provider, callbackURL, w)
	if err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}
	c.Redirect(w, r, resp.Url(rt), resp.Code(http.StatusTemporaryRedirect))
}

func (c *Controller) oauthCallback(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	callbackURL := "/auth/" + provider + "/callback"
	data, err := c.Auth.ExchangeCode(provider, callbackURL, r)
	if err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	user, err := c.User.HandleCallback(data)
	if err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	s, err := c.session(r.Context())
	if err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	if err := s.RegisterUser(w, r, user.ID); err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	c.Redirect(w, r, resp.Success("Welcome, nice to see you here!"), resp.Url(user.HomePath()))
}
