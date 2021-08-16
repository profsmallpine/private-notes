package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
)

type createCommentReq struct {
	Content string `schema:"content,required"`
}

func (c *Controller) createComment(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create comment in the requested group
	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	noteID := mux.Vars(r)[routes.MuxIDParam]
	groupID := mux.Vars(r)[routes.MuxGroupParam]
	rt := fmt.Sprintf("/groups/%s/notes/%s", groupID, noteID)

	note := &domain.Note{}
	if err := c.DB.First(note, noteID).Error; err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	if !user.CanAccessGroup(note.GroupID) {
		err := domain.ErrUnauthorized
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	// Parse + decode form into go
	var req createCommentReq
	if err := c.parseForm(r, &req); err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Setup note for db insert
	comment := &domain.Comment{
		Content: req.Content,
		NoteID:  note.ID,
		UserID:  user.ID,
	}
	if err := c.DB.Create(comment).Error; err != nil {
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	c.Redirect(w, r, resp.Success("Your comment has been successfully created!"), resp.Url(rt))
}
