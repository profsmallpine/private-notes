package web

import (
	"net/http"

	"github.com/xy-planning-network/trails/logger"
)

func (h *Controller) serveWS(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Ranger.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		return
	}

	if err := h.Websocket.Serve(w, r); err != nil {
		h.Ranger.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		return
	}
}
