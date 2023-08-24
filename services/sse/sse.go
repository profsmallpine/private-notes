package sse

import (
	"fmt"
	"net/http"

	"github.com/r3labs/sse/v2"
)

type Service struct {
	*sse.Server
}

func NewService() *Service {
	server := sse.New()
	server.CreateStream("meeting1")

	return &Service{Server: server}
}

func (s *Service) Publish(stream string, data []byte) {
	fmt.Println("publishing")
	s.Server.Publish(stream, &sse.Event{Data: data})
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Server.ServeHTTP(w, r)
}
