package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

type createMeetingReviewReq struct {
	Goals []uint `schema:"goals,required"`
}

func (h *Controller) createMeetingReview(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	rt := fmt.Sprintf("/meetings/%s/review", meetingID)

	meeting := &domain.Meeting{}
	if err := h.DB.First(meeting, meetingID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	_, err = h.User.CanAccessGroup(user, strconv.FormatUint(uint64(meeting.GroupID), 10))
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	// Parse + decode form into go
	var req createMeetingReviewReq
	if err := h.parseForm(r, &req); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	if err := h.Meeting.CopyGoals(meeting, user, req.Goals); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	rt = fmt.Sprintf("/groups/%d/meetings/%s", meeting.GroupID, meetingID)
	h.Redirect(w, r, resp.Url(rt))
}

func (h *Controller) getMeetingReview(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	meetingID := mux.Vars(r)[routes.MuxIDParam]
	rt := fmt.Sprintf("/meetings/%s/review", meetingID)

	meeting := &domain.Meeting{}
	if err := h.DB.Preload("Goals.User").Where("id < ?", meetingID).Order("id DESC").First(meeting).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group, err := h.User.CanAccessGroup(user, strconv.FormatUint(uint64(meeting.GroupID), 10))
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			rt = user.HomePath()
		}
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	data := map[string]any{
		"currentUser":     user,
		"group":           group,
		"meetingID":       meetingID,
		"meetingToReview": meeting,
		"moods":           domain.GoalMoods,
		"styles":          domain.GoalStyles,
	}
	h.Html(
		w,
		r,
		resp.Authed(),
		resp.Data(data),
		resp.Tmpls(
			"tmpl/reviews/show.tmpl",
			"tmpl/goals/_goal.tmpl",
			"tmpl/partials/_header.tmpl",
		),
	)
}
