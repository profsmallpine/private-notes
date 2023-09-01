package web

import (
	"net/http"

	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

func (h *Controller) getRoot(w http.ResponseWriter, r *http.Request) {
	h.Redirect(w, r, resp.Url(routes.GetLoginURL))
}
