package domain

import "net/http"

type WebsocketService interface {
	Broadcast(bytes []byte)
	Serve(w http.ResponseWriter, r *http.Request) error
}
