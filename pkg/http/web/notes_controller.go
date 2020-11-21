package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/http/routes"
)

type createNoteReq struct {
	Content string `schema:"content,required"`
	Title   string `schema:"title,required"`
}

func (c *Controller) createNote(w http.ResponseWriter, r *http.Request) {
	// Authorize user to create note in the requested group
	user := &domain.User{}
	if err := c.DB.Where("id = ?", c.GetCurrentUserID(r)).Preload("Groups").First(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	group := &domain.Group{}
	if err := c.DB.First(group, mux.Vars(r)[routes.MuxIDParam]).Error; err != nil {
		fmt.Println(err)
		return
	}

	if !user.CanAccessGroup(group.ID) {
		err := domain.ErrUnauthorized
		fmt.Println(err)
		return
	}

	// Parse + decode form into go
	var req createNoteReq
	if err := c.parseForm(r, &req); err != nil {
		return
	}

	// Setup note for db insert
	note := &domain.Note{
		Content: req.Content,
		GroupID: group.ID,
		Title:   req.Title,
		UserID:  user.ID,
	}
	if err := c.DB.Create(note).Error; err != nil {
		fmt.Println(err)
		return
	}

	redirect := fmt.Sprintf("/groups/%d/notes/%d", group.ID, note.ID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (c *Controller) getNote(w http.ResponseWriter, r *http.Request) {
	user := &domain.User{}
	if err := c.DB.Preload("Groups").First(user, c.GetCurrentUserID(r)).Error; err != nil {
		fmt.Println(err)
		return
	}

	noteID := mux.Vars(r)[routes.MuxIDParam]
	note := &domain.Note{}
	if err := c.DB.Preload("Comments.Author").Preload("Author").First(note, noteID).Error; err != nil {
		fmt.Println(err)
		return
	}

	// Authorize user can access the requested group
	if !user.CanAccessGroup(note.GroupID) {
		err := domain.ErrUnauthorized
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{"note": note}
	c.respondAuthenticated(w, r, "notes/show", data)
}

func (c *Controller) getNotes(w http.ResponseWriter, r *http.Request) {
	// Authorize user can access the requested group
	user := &domain.User{}
	if err := c.DB.Preload("Groups").First(user, c.GetCurrentUserID(r)).Error; err != nil {
		fmt.Println(err)
		return
	}

	group := &domain.Group{}
	if err := c.DB.First(group, mux.Vars(r)[routes.MuxIDParam]).Error; err != nil {
		fmt.Println(err)
		return
	}

	if !user.CanAccessGroup(group.ID) {
		err := domain.ErrUnauthorized
		fmt.Println(err)
		return
	}

	notes := []*domain.Note{}
	if err := c.DB.Where("group_id = ?", group.ID).Preload("Author").Find(&notes).Error; err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{
		"groupID": group.ID,
		"notes":   notes,
	}
	c.respondAuthenticated(w, r, "notes/index", data)
}

func (c *Controller) newNote(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"groupID": mux.Vars(r)[routes.MuxIDParam]}
	c.respondAuthenticated(w, r, "notes/new", data)
}
