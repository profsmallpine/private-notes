package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

func (h *Controller) oauthLogin(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	callbackURL := "/auth/" + provider + "/callback"
	rt, err := h.Auth.FetchAuthURL(provider, callbackURL, w)
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}
	w.Header().Add("HX-Redirect", rt)
}

func (h *Controller) oauthCallback(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	callbackURL := "/auth/" + provider + "/callback"
	data, err := h.Auth.ExchangeCode(provider, callbackURL, r)
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	user, err := h.User.HandleCallback(data)
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	s, err := h.session(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	if err := s.RegisterUser(w, r, user.ID); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Code(http.StatusSeeOther), resp.Url(routes.GetLoginURL))
		return
	}

	h.Redirect(w, r, resp.Success("Welcome, nice to see you here!"), resp.Url(user.HomePath()))
}
