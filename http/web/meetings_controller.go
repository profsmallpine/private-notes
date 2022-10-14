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

func (h *Controller) createMeeting(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create meeting in the requested group
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	groupID := mux.Vars(r)[routes.MuxIDParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	if !user.IsAdmin {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group, err := h.User.CanAccessGroup(user, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup meeting for db insert
	meeting := &domain.Meeting{GroupID: group.ID}
	if err := h.DB.Create(meeting).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	redirect := fmt.Sprintf("/groups/%d/meetings/%d", group.ID, meeting.ID)
	h.Redirect(w, r, resp.Success("Your meeting has been successfully created!"), resp.Url(redirect))
}

func (h *Controller) getMeeting(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	meeting := &domain.Meeting{}
	if err := h.DB.Preload("Goals.User").First(meeting, meetingID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group, err := h.User.CanAccessGroup(user, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	data := map[string]interface{}{
		"currentUser": user,
		"group":       group,
		"meeting":     meeting,
		"moods":       domain.GoalMoods,
		"styles":      domain.GoalStyles,
	}
	h.Html(
		w,
		r,
		resp.Authed(),
		resp.Data(data),
		resp.Tmpls(
			"tmpl/meetings/show.tmpl",
			"tmpl/goals/_goal.tmpl",
			"tmpl/partials/_header.tmpl",
		),
	)
}

func (h *Controller) updateMeeting(w http.ResponseWriter, r *http.Request) {
	// Authorize user to update meeting in the requested group
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s/meetings/%s", groupID, meetingID)

	if !user.IsAdmin {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	meeting := &domain.Meeting{}
	if err := h.DB.First(meeting, meetingID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt), resp.Code(http.StatusInternalServerError))
		return
	}

	_, err = h.User.CanAccessGroup(user, groupID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup meeting for db udpate
	meeting.Status = domain.MeetingComplete
	if err := h.DB.Save(meeting).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	h.Redirect(w, r, resp.Success("Your meeting has been successfully completed!"), resp.Url(rt))
}
