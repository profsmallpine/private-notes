package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

// Set a Decoder instance as a package global, because it caches meta-data about structs, and an
// instance can be shared safely.
var decoder = schema.NewDecoder()

func (h *Controller) parseForm(r *http.Request, reqStructPtr any) error {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := decoder.Decode(reqStructPtr, r.PostForm); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
