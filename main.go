package main

import (
	"log"
	"os"

	"github.com/profsmallpine/private-notes/app"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	privateNotes, err := app.New(logger)
	if err != nil {
		logger.Fatal(err)
	}

	if err := privateNotes.ListenAndServe(); err != nil {
		logger.Fatal("could not listen and serve: ", err)
	}
}
