package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/pkg/http/middleware"
	"github.com/profsmallpine/private-notes/pkg/http/routes"
)

type route struct {
	path        string
	method      string
	handler     http.HandlerFunc
	middlewares []mux.MiddlewareFunc
}

func (c *Controller) RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	forEveryRequest := []mux.MiddlewareFunc{
		// middleware.LogRequest(c.Logger),
		c.LoadCurrentUserID(),
	}

	// Register unauthenticated routes
	requireUnauthenticatedUser := append(forEveryRequest, c.RequireUnauthenticatedUser())
	unauthenticatedRoutes := []route{
		{path: routes.AuthCallbackURL, method: http.MethodGet, handler: c.oauthCallback},
		{path: routes.AuthLoginURL, method: http.MethodGet, handler: c.oauthLogin},
		{path: routes.GetLoginURL, method: http.MethodGet, handler: c.getLogin},
	}

	for _, r := range unauthenticatedRoutes {
		mws := append(requireUnauthenticatedUser, r.middlewares...)
		router.Handle(r.path, middleware.Chain(
			r.handler,
			mws...,
		)).Methods(r.method)
	}

	// Register authenticated routes
	requireAuthenticatedUser := append(forEveryRequest, c.RequireAuthenticatedUser())
	authenticatedRoutes := []route{
		{path: routes.CreateCommentURL, method: http.MethodPost, handler: c.createComment},
		{path: routes.CreateGroupURL, method: http.MethodPost, handler: c.createGroup},
		{path: routes.CreateNoteURL, method: http.MethodPost, handler: c.createNote},
		{path: routes.GetGroupsURL, method: http.MethodGet, handler: c.getGroups},
		{path: routes.GetLogoffURL, method: http.MethodGet, handler: c.getLogoff},
		{path: routes.GetNoteURL, method: http.MethodGet, handler: c.getNote},
		{path: routes.GetNotesURL, method: http.MethodGet, handler: c.getNotes},
		{path: routes.GetRootURL, method: http.MethodGet, handler: c.getRoot},
		{path: routes.NewGroupURL, method: http.MethodGet, handler: c.newGroup},
		{path: routes.NewNoteURL, method: http.MethodGet, handler: c.newNote},
	}

	for _, r := range authenticatedRoutes {
		mws := append(requireAuthenticatedUser, r.middlewares...)
		router.Handle(r.path, middleware.Chain(
			r.handler,
			mws...,
		)).Methods(r.method)
	}

	return router
}
