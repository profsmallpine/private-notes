package domain

import "net/http"

type SSEService interface {
	Publish(stream string, bytes []byte)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
