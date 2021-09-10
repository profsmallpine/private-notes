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
	"github.com/profsmallpine/private-notes/services/email"
	"github.com/xy-planning-network/trails/http/router"
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/logger"
	"github.com/xy-planning-network/trails/postgres"
)

type App struct {
	*http.Server
}

func New(logging *log.Logger) (*App, error) {
	_ = godotenv.Load()

	allowedEmails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")
	baseURL := envVarOrString("BASE_URL", "http://localhost:8080")
	env := os.Getenv("ENVIRONMENT")
	port := envVarOrString("PORT", ":8080")
	if port[0] != ':' {
		port = ":" + port
	}

	// Connect/migrate database.
	config := &postgres.CxnConfig{IsTestDB: false, URL: os.Getenv("DATABASE_URL")}
	if config.URL == "" {
		config.Host = envVarOrString("PG_HOST", "localhost")
		config.Name = envVarOrString("PG_NAME", "private_notes_dev_db")
		config.Password = envVarOrString("PG_PASSWORD", "")
		config.Port = envVarOrString("PG_PORT", "5432")
		config.User = envVarOrString("PG_USER", "private_notes_dev_db_user")
	}
	db, err := postgres.Connect(config, migrations.List)
	if err != nil {
		return nil, err
	}

	sss, err := session.NewStoreService(env, os.Getenv("SESSION_AUTH_KEY"), os.Getenv("SESSION_ENCRYPTION_KEY"))
	if err != nil {
		return nil, err
	}

	es := email.NewService(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_PORT"),
	)

	ls := logger.NewLogger(
		logger.WithEnv(env),
		logger.WithLogger(logging),
		logger.WithLevel(logger.LogLevelInfo),
	)

	services := domain.Services{
		Auth:         auth.NewService(baseURL),
		Email:        es,
		Logger:       ls,
		SessionStore: sss,
	}

	procedures := procedures.New(allowedEmails, db, services)

	controller := web.Controller{DB: db, Procedures: procedures, Services: services}

	r := controller.Router(env, baseURL)

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
