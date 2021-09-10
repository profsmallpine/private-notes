package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/logger"
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
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	group := &domain.Group{}
	if err := c.DB.Preload("Users").First(group, groupID).Error; err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	if !user.CanAccessGroup(note.GroupID) {
		err := domain.ErrUnauthorized
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(user.HomePath()))
		return
	}

	// Parse + decode form into go
	var req createCommentReq
	if err := c.parseForm(r, &req); err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
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
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(rt))
		return
	}

	// Send email to other members of group
	members := []string{}
	for _, u := range group.Users {
		if u.ID != user.ID {
			members = append(members, u.Email)
		}
	}
	if len(members) > 0 {
		// TODO: log failure
		msg := fmt.Sprintf("%s commented on a note. Check it out here: %s", user.FirstName, os.Getenv("BASE_URL")+rt)
		subject := "Time to Reflect! Comment Made on Note"
		if err := c.Services.Email.Send(msg, subject, members); err != nil {
			c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		}
	}

	c.Redirect(w, r, resp.Success("Your comment has been successfully created!"), resp.Url(rt))
}
