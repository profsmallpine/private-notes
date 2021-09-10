package web

import (
	"fmt"
	"net/http"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/http/routes"
	"github.com/xy-planning-network/trails/http/resp"
	"github.com/xy-planning-network/trails/logger"
)

type createGroupReq struct {
	Description string `schema:"description,required"`
	Name        string `schema:"name,required"`
	UserIDs     []uint `schema:"userIDs,required"`
}

func (c *Controller) createGroup(w http.ResponseWriter, r *http.Request) {
	// Parse + decode form into go
	var req createGroupReq
	if err := c.parseForm(r, &req); err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetGroupsURL))
		return
	}

	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
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
	if err := c.DB.Create(group).Error; err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetGroupsURL))
		return
	}

	rt := fmt.Sprintf("%s/%d/notes", routes.GetGroupsURL, group.ID)
	c.Redirect(w, r, resp.Success("You've successfully created a new group, woo hoo!"), resp.Url(rt))
}

func (c *Controller) getGroups(w http.ResponseWriter, r *http.Request) {
	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	groups := []*domain.Group{}
	if err := c.DB.Model(user).Preload("Users").Association("Groups").Find(&groups); err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetLogoffURL))
		return
	}

	data := map[string]interface{}{"groups": groups}
	c.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls("tmpl/groups/index.tmpl", "tmpl/partials/_header.tmpl"))
}

func (c *Controller) newGroup(w http.ResponseWriter, r *http.Request) {
	user, err := c.currentUser(r.Context())
	if err != nil {
		c.Redirect(w, r, resp.Url(routes.GetLogoffURL))
		return
	}

	users := []*domain.User{}
	if err := c.DB.Where("id != ?", user.ID).Find(&users).Error; err != nil {
		c.Logger.Error(err.Error(), &logger.LogContext{Request: r, User: user, Error: err})
		c.Redirect(w, r, resp.GenericErr(err), resp.Url(routes.GetGroupsURL))
		return
	}

	data := map[string]interface{}{"users": users}
	c.Html(w, r, resp.Authed(), resp.Data(data), resp.Tmpls("tmpl/groups/new.tmpl", "tmpl/partials/_header.tmpl"))
}
