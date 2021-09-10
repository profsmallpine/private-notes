package web

import (
	"net/http"
	"net/url"

	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/ctx"
	"github.com/xy-planning-network/trails/http/middleware"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/http/router"
	"github.com/xy-planning-network/trails/http/template"
)

func (c *Controller) Router(env, baseURL string) router.Router {
	c.Keyring = ctx.NewKeyRing(sessionCtxKey, currentUserCtxKey)

	r := router.NewRouter(env)

	u, _ := url.ParseRequestURI(baseURL)
	p := template.NewParser(
		files,
		template.WithFn(template.Env(env)),
		template.WithFn(template.RootUrl(u)),
		// template.WithFn(template.Nonce()),
		// template.WithFn("packTag", template.TagPacker(h.Env.String(), os.DirFS("."))),
		// template.WithFn("isDevelopment", func() bool { return h.Env == domain.Development }),
		// template.WithFn("isStaging", func() bool { return h.Env == domain.Staging }),
		// template.WithFn("isProduction", func() bool { return h.Env == domain.Production }),
	)

	c.Responder = resp.NewResponder(
		resp.WithRootUrl(baseURL),
		resp.WithLogger(c.Logger),
		resp.WithParser(p),
		resp.WithAuthTemplate(authedTmpl),
		resp.WithUnauthTemplate(unauthedTmpl),
		resp.WithSessionKey(c.Keyring.SessionKey()),
		resp.WithUserSessionKey(c.Keyring.CurrentUserKey()),
		resp.WithContactErrMsg(errorMessage),
	)

	r.OnEveryRequest(
		middleware.ForceHTTPS(env),
		middleware.InjectIPAddress(),
		middleware.LogRequest(c.Logger),
		middleware.InjectSession(c.SessionStore, c.Keyring.SessionKey()),
		middleware.CurrentUser(c.Responder, c, c.Keyring.SessionKey(), c.Keyring.CurrentUserKey()),
	)

	// Register unauthenticated routes
	unauthenticatedRoutes := []router.Route{
		{Path: routes.AuthCallbackURL, Method: http.MethodGet, Handler: c.oauthCallback},
		{Path: routes.AuthLoginURL, Method: http.MethodGet, Handler: c.oauthLogin},
		{Path: routes.GetLoginURL, Method: http.MethodGet, Handler: c.getLogin},
	}
	r.UnauthedRoutes(c.Keyring.CurrentUserKey(), unauthenticatedRoutes)

	// Register authenticated routes
	authenticatedRoutes := []router.Route{
		{Path: routes.CreateCommentURL, Method: http.MethodPost, Handler: c.createComment},
		{Path: routes.CreateGroupURL, Method: http.MethodPost, Handler: c.createGroup},
		{Path: routes.CreateNoteURL, Method: http.MethodPost, Handler: c.createNote},
		{Path: routes.GetGroupsURL, Method: http.MethodGet, Handler: c.getGroups},
		{Path: routes.GetLogoffURL, Method: http.MethodGet, Handler: c.getLogoff},
		{Path: routes.GetNoteURL, Method: http.MethodGet, Handler: c.getNote},
		{Path: routes.GetNotesURL, Method: http.MethodGet, Handler: c.getNotes},
		{Path: routes.GetRootURL, Method: http.MethodGet, Handler: c.getRoot},
		{Path: routes.NewGroupURL, Method: http.MethodGet, Handler: c.newGroup},
		{Path: routes.NewNoteURL, Method: http.MethodGet, Handler: c.newNote},
	}
	r.AuthedRoutes(c.Keyring.CurrentUserKey(), routes.GetLoginURL, routes.GetLogoffURL, authenticatedRoutes)

	return r
}
