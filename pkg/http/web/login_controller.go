package web

import (
	"net/http"

	"github.com/profsmallpine/private-notes/pkg/http/routes"
)

func (c *Controller) getLogin(w http.ResponseWriter, r *http.Request) {
	c.respondUnauthenticated(w, r, "login/index", nil)
}

func (c *Controller) getLogoff(w http.ResponseWriter, r *http.Request) {
	c.DeleteSession(w, r)
	http.Redirect(w, r, routes.GetLoginURL, http.StatusFound)
}
