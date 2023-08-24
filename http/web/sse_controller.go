package web

import (
	"net/http"

	"github.com/xy-planning-network/trails/logger"
)

func (h *Controller) serveSSE(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Ranger.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		return
	}

	go func() {
		// Received Browser Disconnection
		<-r.Context().Done()
		println("The client is disconnected")
	}()

	h.SSE.ServeHTTP(w, r)
}
