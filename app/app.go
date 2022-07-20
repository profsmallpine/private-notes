package app

import (
	"embed"
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
	"github.com/xy-planning-network/trails/http/session"
	"github.com/xy-planning-network/trails/http/template"
	"github.com/xy-planning-network/trails/postgres"
	"github.com/xy-planning-network/trails/ranger"
)

func New(logging *log.Logger, files embed.FS) (*ranger.Ranger, error) {
	_ = godotenv.Load()

	allowedEmails := strings.Split(os.Getenv("ALLOWED_EMAILS"), ",")
	baseURL := envVarOrString("BASE_URL", "http://localhost:8080")
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

	es := email.NewService(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_PORT"),
	)

	services := domain.Services{
		Auth:  auth.NewService(baseURL),
		Email: es,
	}

	procedures := procedures.New(allowedEmails, db, services)

	// Configure redis for sessions
	redisURL := os.Getenv("REDIS_URL")
	redisPassword := ""
	redisURI := "localhost:6379"
	if redisURL != "" {
		parts := strings.Split(redisURL, "@")
		redisPassword = strings.Split(parts[0], ":")[2]
		redisURI = parts[1]
	}
	fsOpt := session.WithRedis(redisURI, redisPassword)

	// Configure http server
	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Configure ranger
	rng, err := ranger.New(
		ranger.DefaultSessionStore(fsOpt),
		ranger.DefaultParser(
			files,
			template.WithFn("add1", func(value int) int { return value + 1 }),
			template.WithFn("minus1", func(value int) int { return value - 1 }),
		),
		ranger.WithDB(postgres.NewService(db)),
		ranger.WithServer(srv),
		ranger.WithUserSessions(procedures.User),
	)
	if err != nil {
		return nil, err
	}

	services.Logger = rng.Logger

	h := &web.Controller{
		DB:         db,
		Procedures: procedures,
		Ranger:     rng,
		Services:   services,
	}

	h.Router()

	return rng, nil
}

func envVarOrString(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
