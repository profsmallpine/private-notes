package websocket

import "net/http"

type Service struct {
	*Hub
}

func NewService() *Service {
	hub := newHub()
	go hub.run()
	return &Service{Hub: hub}
}

func (s *Service) Broadcast(bytes []byte) {
	s.Hub.broadcast <- bytes
}

// Serve handles websocket requests from the peer.
func (s *Service) Serve(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{hub: s.Hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	return nil
}
