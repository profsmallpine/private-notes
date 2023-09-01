package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

type createGroupReq struct {
	Description string `schema:"description,required"`
	Name        string `schema:"name,required"`
	UserIDs     []uint `schema:"userIDs,required"`
}

func (h *Controller) createGroup(w http.ResponseWriter, r *http.Request) {
	// Parse + decode form into go
	var req createGroupReq
	if err := h.parseForm(r, &req); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.NewGroupURL), resp.Code(http.StatusBadRequest))
		return
	}

	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL), resp.Code(http.StatusUnauthorized))
		return
	}

	// NOTE: passing in the current user here so it includes the current user as a member of the group.
	users := []*domain.User{user}
	for _, id := range req.UserIDs {
		users = append(users, domain.NewUserFromID(id))
	}
	group := &domain.Group{
		Description: req.Description,
		Name:        req.Name,
		Users:       users,
	}
	if err := h.DB.Create(group).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.NewGroupURL), resp.Code(http.StatusInternalServerError))
		return
	}

	rt := fmt.Sprintf("%s/%d", routes.GetGroupsURL, group.ID)
	h.Redirect(w, r, resp.Success("You've successfully created a new group, woo hoo!"), resp.Url(rt))
}

func (h *Controller) getGroup(w http.ResponseWriter, r *http.Request) {
	// Authorize user can access the requested group
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	group := &domain.Group{}
	if err := h.DB.First(group, mux.Vars(r)[routes.MuxIDParam]).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	if !user.CanAccessGroup(group.ID) {
		err := domain.ErrUnauthorized
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	query := "group_id = ?"
	params := []any{group.ID}
	order := "created_at DESC"

	meetings := []*domain.Meeting{}
	meetingsPD, err := h.Ranger.DB().PagedByQuery(&meetings, query, params, order, 1, domain.PerPageSize)
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	notes := []*domain.Note{}
	notesPD, err := h.Ranger.DB().PagedByQuery(&notes, query, params, order, 1, domain.PerPageSize, "Author")
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	data := map[string]any{
		"currentUser": user,
		"groupID":     group.ID,
		"meetings":    meetingsPD,
		"notes":       notesPD,
	}
	h.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls(
		"tmpl/groups/show.tmpl",
		"tmpl/partials/_header.tmpl",
		"tmpl/notes/_list.tmpl",
		"tmpl/meetings/_list.tmpl",
	))
}

func (h *Controller) getGroups(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	groups := []*domain.Group{}
	if err := h.DB.Model(user).Preload("Users").Association("Groups").Find(&groups); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetLogoffURL))
		return
	}

	data := map[string]any{"groups": groups}
	h.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls("tmpl/groups/index.tmpl", "tmpl/partials/_header.tmpl"))
}

func (h *Controller) newGroup(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	users := []*domain.User{}
	if err := h.DB.Where("id != ?", user.ID).Find(&users).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetGroupsURL))
		return
	}

	data := map[string]any{"users": users}
	h.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls("tmpl/groups/new.tmpl", "tmpl/partials/_header.tmpl"))
}
