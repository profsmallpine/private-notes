package app

import (
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/auth"
	"github.com/profsmallpine/private-notes/pkg/http/web"
	"github.com/profsmallpine/private-notes/pkg/postgres"
	"github.com/profsmallpine/private-notes/pkg/procedures"
)

type App struct {
	*http.Server
}

func New(logger *log.Logger) (*App, error) {
	_ = godotenv.Load()

	allowedEmails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")

	baseURL := envVarOrString("PORT", "http://localhost:8080")
	port := envVarOrString("PORT", ":8080")
	pgConnStr := envVarOrString(
		"DATABASE_URL",
		"host=localhost port=5432 user=private_notes_dev_db_user dbname=private_notes_dev_db sslmode=disable",
	)

	authKey, err := hex.DecodeString(os.Getenv("SESSION_AUTH_KEY"))
	if err != nil {
		return nil, err
	}
	encryptionKey, err := hex.DecodeString(os.Getenv("SESSION_ENCRYPTION_KEY"))
	if err != nil {
		return nil, err
	}

	db, err := postgres.ConnectDatabase(pgConnStr)
	if err != nil {
		return nil, err
	}

	services := domain.Services{
		Auth: auth.NewService(baseURL),
	}

	procedures := procedures.New(allowedEmails, db, services)

	controller := web.NewController(
		db,
		procedures,
		services,
		authKey,
		encryptionKey,
	)

	router := controller.RegisterRoutes()

	server := newServer(port, router)

	return &App{server}, nil
}

func newServer(port string, router *mux.Router) *http.Server {
	return &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func envVarOrString(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
