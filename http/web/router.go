package web

import (
	"net/http"

	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/router"
)

func (h *Controller) Router() {
	unauthenticatedRoutes := []router.Route{
		{Path: routes.AuthCallbackURL, Method: http.MethodGet, Handler: h.oauthCallback},
		{Path: routes.AuthLoginURL, Method: http.MethodGet, Handler: h.oauthLogin},
		{Path: routes.GetLoginURL, Method: http.MethodGet, Handler: h.getLogin},
	}
	h.UnauthedRoutes(h.EmitKeyring().CurrentUserKey(), unauthenticatedRoutes)

	// Register authenticated routes
	authenticatedRoutes := []router.Route{
		{Path: routes.CreateCommentURL, Method: http.MethodPost, Handler: h.createComment},
		{Path: routes.CreateGoalURL, Method: http.MethodPost, Handler: h.createGoal},
		{Path: routes.CreateGroupURL, Method: http.MethodPost, Handler: h.createGroup},
		{Path: routes.CreateMeetingURL, Method: http.MethodGet, Handler: h.createMeeting}, // Kind of hacky, but no data is passed to this endpoint yet
		{Path: routes.CreateNoteURL, Method: http.MethodPost, Handler: h.createNote},
		{Path: routes.GetGroupURL, Method: http.MethodGet, Handler: h.getGroup},
		{Path: routes.GetGroupsURL, Method: http.MethodGet, Handler: h.getGroups},
		{Path: routes.GetLogoffURL, Method: http.MethodGet, Handler: h.getLogoff},
		{Path: routes.GetNoteURL, Method: http.MethodGet, Handler: h.getNote},
		{Path: routes.GetMeetingURL, Method: http.MethodGet, Handler: h.getMeeting},
		{Path: routes.GetNotesURL, Method: http.MethodGet, Handler: h.getNotes},
		{Path: routes.GetRootURL, Method: http.MethodGet, Handler: h.getRoot},
		{Path: routes.NewGroupURL, Method: http.MethodGet, Handler: h.newGroup},
		{Path: routes.NewNoteURL, Method: http.MethodGet, Handler: h.newNote},
	}
	h.AuthedRoutes(h.EmitKeyring().CurrentUserKey(), routes.GetLoginURL, routes.GetLogoffURL, authenticatedRoutes)
}
