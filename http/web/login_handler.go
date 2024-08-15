package web

import (
	"net/http"

	"github.com/profsmallpine/private-notes/html"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

func (h *Controller) getLogin(w http.ResponseWriter, r *http.Request) {
	html.UnauthenticatedLayout(
		h.flashes(w, r),
		html.Login(),
	).Render(r.Context(), w)
}

func (h *Controller) getLogoff(w http.ResponseWriter, r *http.Request) {
	s, err := h.session(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Err(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	if err := s.Delete(w, r); err != nil {
		h.Redirect(w, r, resp.Err(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	h.Redirect(w, r, resp.Url(routes.GetLoginURL))
}
