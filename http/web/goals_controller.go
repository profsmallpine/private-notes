package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

type createGoalReq struct {
	Content string `schema:"content,required"`
	Mood    string `schema:"mood,required"`
}

func (h *Controller) createGoal(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create goal in the requested meeting
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s/meetings/%s", groupID, meetingID)

	meeting := &domain.Meeting{}
	if err := h.DB.First(meeting, meetingID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group := &domain.Group{}
	if err := h.DB.Preload("Users").First(group, groupID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	if !user.CanAccessGroup(meeting.GroupID) {
		err := domain.ErrUnauthorized
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	// Parse + decode form into go
	var req createGoalReq
	if err := h.parseForm(r, &req); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup goal for db insert
	goal := &domain.Goal{
		Content:   req.Content,
		MeetingID: meeting.ID,
		Mood:      req.Mood,
		UserID:    user.ID,
	}
	if err := h.DB.Create(goal).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	h.Redirect(w, r, resp.Success("Your goal has been successfully created!"), resp.Url(rt))
}
