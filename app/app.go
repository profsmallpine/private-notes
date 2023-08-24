package app

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/web"
	"github.com/profsmallpine/private-notes/migrations"
	"github.com/profsmallpine/private-notes/procedures"
	"github.com/profsmallpine/private-notes/services/auth"
	"github.com/profsmallpine/private-notes/services/email"
	"github.com/profsmallpine/private-notes/services/websocket"
	"github.com/xy-planning-network/trails"
	"github.com/xy-planning-network/trails/postgres"
	"github.com/xy-planning-network/trails/ranger"
)

func New(logging *log.Logger, files embed.FS) (*ranger.Ranger, error) {
	isMaintMode := os.Getenv("MAINTENANCE_MODE") == "true"

	cfg := ranger.Config[*domain.User]{
		FS:         files,
		MaintMode:  isMaintMode,
		Migrations: migrations.List,
	}

	rng, err := ranger.New(cfg)
	if err != nil {
		return nil, err
	}

	pdb, ok := rng.DB().(*postgres.DatabaseServiceImpl)
	if !ok {
		err := fmt.Errorf("%w: rng.DB() is not *postgres.DatabaseServiceImpl, but is %T", trails.ErrBadConfig, rng.DB())
		return nil, err
	}

	es := email.NewService(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_PORT"),
	)

	services := domain.Services{
		Auth:  auth.NewService(trails.EnvVarOrString("BASE_URL", "http://localhost:8080")),
		Email: es,
		// SSE:       sse.NewService(),
		Websocket: websocket.NewService(),
	}

	procedures := procedures.New(strings.Split(os.Getenv("ALLOWED_EMAILS"), ","), pdb.DB, services)

	services.Logger = rng.Logger

	h := &web.Controller{
		DB:         pdb.DB,
		Procedures: procedures,
		Ranger:     rng,
		Services:   services,
	}

	h.Router()

	return rng, nil
}
