package web

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/html"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/logger"
)

type createNoteReq struct {
	Content string `schema:"content,required"`
	Title   string `schema:"title,required"`
}

func (h *Controller) createNote(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create note in the requested group
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	group := &domain.Group{}
	groupID := mux.Vars(r)[routes.MuxIDParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	if err := h.DB.Preload("Users").First(group, groupID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	if !user.CanAccessGroup(group.ID) {
		err := domain.ErrUnauthorized
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	// Parse + decode form into go
	var req createNoteReq
	if err := h.parseForm(r, &req); err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup note for db insert
	note := &domain.Note{
		Content: req.Content,
		GroupID: group.ID,
		Title:   req.Title,
		UserID:  user.ID,
	}
	if err := h.DB.Create(note).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Send email to other members of group
	redirect := fmt.Sprintf("/groups/%d/notes/%d", group.ID, note.ID)
	members := []string{}
	for _, u := range group.Users {
		if u.ID != user.ID {
			members = append(members, u.Email)
		}
	}
	if len(members) > 0 {
		msg := fmt.Sprintf("%s created a new post. Check it out here: %s", user.FirstName, os.Getenv("BASE_URL")+redirect)
		subject := fmt.Sprintf("Get Excited! Newness Submitted to %s", group.Name)
		if err := h.Services.Email.Send(msg, subject, members); err != nil {
			h.Ranger.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		}
	}

	h.Redirect(w, r, resp.Success("Your note has been successfully created!"), resp.Url(redirect))
}

func (h *Controller) getNote(w http.ResponseWriter, r *http.Request) {
	user, err := h.currentUser(r.Context())
	if err != nil {
		h.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	noteID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s", groupID)

	note := &domain.Note{}
	if err := h.DB.Preload("Comments.Author").Preload("Author").Preload("Group").First(note, noteID).Error; err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Authorize user can access the requested group
	if !user.CanAccessGroup(note.GroupID) {
		err := domain.ErrUnauthorized
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	html.AuthenticatedLayout(
		h.flashes(w, r),
		html.ShowNote(note),
		[]domain.Breadcrumb{{Label: note.Group.Name, URL: fmt.Sprintf("/groups/%d", note.Group.ID)}},
	).Render(r.Context(), w)
}

func (h *Controller) getNotes(w http.ResponseWriter, r *http.Request) {
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

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page == 0 {
		page = 1
	}

	query := "group_id = ?"
	params := []any{group.ID}
	order := "created_at DESC"

	notes := []*domain.Note{}
	pd, err := h.Ranger.DB().PagedByQuery(&notes, query, params, order, page, domain.PerPageSize, "Author")
	if err != nil {
		h.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	html.ListNotes(pd, group).Render(r.Context(), w)
}

func (h *Controller) newNote(w http.ResponseWriter, r *http.Request) {
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

	html.AuthenticatedLayout(
		h.flashes(w, r),
		html.NewNote(group),
		[]domain.Breadcrumb{{Label: group.Name, URL: fmt.Sprintf("/groups/%d", group.ID)}},
	).Render(r.Context(), w)
}
