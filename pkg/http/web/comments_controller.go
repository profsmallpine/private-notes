package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/http/routes"
)

type createCommentReq struct {
	Content string `schema:"content,required"`
}

func (c *Controller) createComment(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create comment in the requested group
	user := &domain.User{}
	if err := c.DB.Where("id = ?", c.GetCurrentUserID(r)).Preload("Groups").First(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	note := &domain.Note{}
	if err := c.DB.First(note, mux.Vars(r)[routes.MuxIDParam]).Error; err != nil {
		fmt.Println(err)
		return
	}

	if !user.CanAccessGroup(note.GroupID) {
		err := domain.ErrUnauthorized
		fmt.Println(err)
		return
	}

	// Parse + decode form into go
	var req createCommentReq
	if err := c.parseForm(r, &req); err != nil {
		return
	}

	// Setup note for db insert
	comment := &domain.Comment{
		Content: req.Content,
		NoteID:  note.ID,
		UserID:  user.ID,
	}
	if err := c.DB.Create(comment).Error; err != nil {
		fmt.Println(err)
		return
	}

	// TODO: flash for success
	redirect := fmt.Sprintf("/groups/%d/notes/%d", note.GroupID, note.ID)
	http.Redirect(w, r, redirect, http.StatusFound)
}
