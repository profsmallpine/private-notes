package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

func (c *Controller) createMeeting(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create meeting in the requested group
	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	groupID := mux.Vars(r)[routes.MuxIDParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	if !user.IsAdmin {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group, err := c.User.CanAccessGroup(user, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup meeting for db insert
	meeting := &domain.Meeting{GroupID: group.ID}
	if err := c.DB.Create(meeting).Error; err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	redirect := fmt.Sprintf("/groups/%d/meetings/%d", group.ID, meeting.ID)
	c.Redirect(w, r, resp.Success("Your meeting has been successfully created!"), resp.Url(redirect))
}

func (c *Controller) getMeeting(w http.ResponseWriter, r *http.Request) {
	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	meeting := &domain.Meeting{}
	if err := c.DB.Preload("Goals.User").First(meeting, meetingID).Error; err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group, err := c.User.CanAccessGroup(user, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	data := map[string]interface{}{
		"group":   group,
		"meeting": meeting,
		"moods":   domain.GoalMoods,
	}
	c.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls("tmpl/meetings/show.tmpl", "tmpl/partials/_header.tmpl"))
}
