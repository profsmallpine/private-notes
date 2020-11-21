package web

import (
	"fmt"
	"net/http"

	"github.com/profsmallpine/private-notes/domain"
	"github.com/profsmallpine/private-notes/pkg/http/routes"
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
		return
	}

	// Setup group for db insert
	users := []*domain.User{
		// NOTE: passing in the current user ID here so it includes the current user as a member of
		// the group.
		domain.NewUserFromID(c.GetCurrentUserID(r)),
	}
	for _, id := range req.UserIDs {
		users = append(users, domain.NewUserFromID(id))
	}
	group := &domain.Group{
		Description: req.Description,
		Name:        req.Name,
		Users:       users,
	}
	if err := c.DB.Create(group).Error; err != nil {
		fmt.Println(err)
		return
	}

	redirect := fmt.Sprintf("%s/%d", routes.GetGroupsURL, group.ID)
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (c *Controller) getGroups(w http.ResponseWriter, r *http.Request) {
	groups := []*domain.Group{}
	user := domain.NewUserFromID(c.GetCurrentUserID(r))
	if err := c.DB.Model(user).Association("Groups").Find(&groups); err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{"groups": groups}
	c.respondAuthenticated(w, r, "groups/index", data)
}

func (c *Controller) newGroup(w http.ResponseWriter, r *http.Request) {
	users := []*domain.User{}
	if err := c.DB.Where("id != ?", c.GetCurrentUserID(r)).Find(&users).Error; err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{"users": users}
	c.respondAuthenticated(w, r, "groups/new", data)
}
