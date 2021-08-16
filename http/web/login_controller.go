package web

import (
	"net/http"

	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

func (c *Controller) getLogin(w http.ResponseWriter, r *http.Request) {
	c.Html(w, r, resp.Unauthed(), resp.Tmpls("tmpl/login/index.tmpl"))
}

func (c *Controller) getLogoff(w http.ResponseWriter, r *http.Request) {
	s, err := c.session(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Err(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	if err := s.Delete(w, r); err != nil {
		c.Redirect(w, r, resp.Err(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	c.Redirect(w, r, resp.Url(routes.GetLoginURL))
}
