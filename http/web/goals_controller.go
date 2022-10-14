package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

type createGoalReq struct {
	Content string           `schema:"content,required"`
	Mood    string           `schema:"mood,required"`
	Style   domain.GoalStyle `schema:"style,required"`
}

func (h *Controller) createGoal(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create goal in the requested meeting
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL), resp.Code(http.StatusUnauthorized))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s/meetings/%s", groupID, meetingID)

	meeting := &domain.Meeting{}
	if err := h.DB.First(meeting, meetingID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt), resp.Code(http.StatusInternalServerError))
		return
	}

	group := &domain.Group{}
	if err := h.DB.Preload("Users").First(group, groupID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt), resp.Code(http.StatusInternalServerError))
		return
	}

	if !user.CanAccessGroup(meeting.GroupID) {
		err := domain.ErrUnauthorized
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()), resp.Code(http.StatusUnauthorized))
		return
	}

	if meeting.IsComplete() {
		err := domain.ErrBadRequest
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt), resp.Code(http.StatusBadRequest))
		return
	}

	// Parse + decode form into go
	var req createGoalReq
	if err := h.parseForm(r, &req); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Default to goal start style if invalid
	if req.Style.Valid() != nil {
		req.Style = domain.GoalStart
	}

	// Setup goal for db insert
	goal := &domain.Goal{
		Content:   req.Content,
		MeetingID: meeting.ID,
		Mood:      req.Mood,
		Style:     req.Style,
		User:      user,
	}
	if err := h.DB.Create(goal).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt), resp.Code(http.StatusInternalServerError))
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "text/vnd.turbo-stream.html") {
		tmpl, _ := template.ParseFiles("tmpl/goals/goal.turbo.tmpl", "tmpl/goals/_goal.tmpl")
		var buf bytes.Buffer
		tmpl.Execute(&buf, goal)
		h.Websocket.Broadcast(buf.Bytes())
		return
	}

	h.Redirect(w, r, resp.Success("Your goal has been successfully created!"), resp.Url(rt))
}
