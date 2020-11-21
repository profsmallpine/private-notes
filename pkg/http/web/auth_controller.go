package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) oauthLogin(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	callbackURL := "/auth/" + provider + "/callback"
	rt, err := c.Auth.FetchAuthURL(provider, callbackURL, w)
	if err != nil {
		// c.redirectWithGenericError(w, r, err, nil, routes.GetLoginURL)
		return
	}
	http.Redirect(w, r, rt, http.StatusTemporaryRedirect)
}

func (c *Controller) oauthCallback(w http.ResponseWriter, r *http.Request) {
	provider := mux.Vars(r)["id"]

	// 	if r.URL.Query().Get("error") != "" {
	// 		h.redirectWithCustomError(w, r, nil, nil, routes.GetLoginURL, FlashInfo, DeniedAccessMessage)
	// 		return
	// 	}

	callbackURL := "/auth/" + provider + "/callback"
	data, err := c.Auth.ExchangeCode(provider, callbackURL, r)
	if err != nil {
		// // NOTE: it seems that around 1 of every 1000 login requests does not return
		// // a cookie for validation. We can skip the notification to sentry in these
		// // cases.
		// if errors.Is(err, domain.ErrorInvalidAccess) {
		// 	h.redirectWithGenericError(w, r, nil, nil, routes.GetLoginURL)
		// 	return
		// }

		// h.redirectWithGenericError(w, r, err, nil, routes.GetLoginURL)
		return
	}

	user, err := c.User.HandleCallback(data)
	if err != nil {
		return
	}

	c.RegisterSession(w, r, user.ID)

	// 	user, err := h.Procedures.Login.Create(token.ServiceID, middleware.GetIPAddress(r), nil)
	// 	if err != nil {
	// 		if errors.Is(err, domain.ErrorNotFound) {
	// 			h.setFlash(w, r, FlashInfo, NoAccountMessage)
	// 		} else {
	// 			h.setFlash(w, r, FlashError, ErrorMessage)
	// 		}
	// 		http.Redirect(w, r, routes.GetLoginURL, http.StatusSeeOther)
	// 		return
	// 	}

	// 	if err := h.RegisterUserSession(w, r, user.ID); err != nil {
	// 		http.Redirect(w, r, routes.GetLoginURL, http.StatusSeeOther)
	// 		return
	// 	}

	// 	h.redirectOnSuccess(w, r, user.HomePath(), WelcomeMessage)
	http.Redirect(w, r, user.HomePath(), http.StatusFound)
}
