package app

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/web"
	"github.com/profsmallpine/private-notes/migrations"
	"github.com/profsmallpine/private-notes/procedures"
	"github.com/profsmallpine/private-notes/services/auth"
	"github.com/xy-planning-network/trails/http/router"
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/postgres"
)

type App struct {
	*http.Server
}

func New(logger *log.Logger) (*App, error) {
	_ = godotenv.Load()

	allowedEmails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")

	baseURL := envVarOrString("PORT", "http://localhost:8080")
	port := envVarOrString("PORT", ":8080")

	// Connect/migrate database.
	db, err := postgres.Connect(&postgres.CxnConfig{
		Host:     envVarOrString("PG_HOST", "localhost"),
		IsTestDB: false,
		Name:     envVarOrString("PG_NAME", "private_notes_dev_db"),
		Password: envVarOrString("PG_PASSWORD", ""),
		Port:     envVarOrString("PG_PORT", "5432"),
		User:     envVarOrString("PG_USER", "private_notes_dev_db_user"),
	}, migrations.List)
	if err != nil {
		return nil, err
	}

	sss, err := session.NewStoreService("development", os.Getenv("SESSION_AUTH_KEY"), os.Getenv("SESSION_ENCRYPTION_KEY"))
	if err != nil {
		return nil, err
	}

	services := domain.Services{
		Auth:         auth.NewService(baseURL),
		SessionStore: sss,
	}

	procedures := procedures.New(allowedEmails, db, services)

	controller := web.Controller{DB: db, Procedures: procedures, Services: services}

	r := controller.Router()

	server := newServer(port, r)

	return &App{server}, nil
}

// TODO: nicely handle SIGTERM

func newServer(port string, r router.Router) *http.Server {
	return &http.Server{
		Addr:         port,
		Handler:      r,
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
