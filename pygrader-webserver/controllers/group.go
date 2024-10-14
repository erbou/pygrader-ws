package controllers

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"pygrader-webserver/models"
)

type GroupController struct {
	beego.Controller
}

func (c *GroupController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title Create
// @Description create group
// @Param	body		body 	models.Group	true		"The group content"
// @Success 200 {object} models.GroupView
// @Failure 400 invalid input
// @router / [post]
func (c *GroupController) Post() {
	if g, s, err := models.Verify[models.Group](c.Ctx.Input.RequestBody); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if obj, err := models.AddGroup(g); err != nil {
		CustomAbort(c, err, 400, "Failed")
	} else {
		if _, err := models.GroupAddUser(g, s.Issuer, "owner"); err != nil {
			logs.Error("Could not add creator %v to group %v", s.Issuer, g)
		}

		obj, _ := models.GetGroup(obj.Id)
		c.SetData(obj.View())
	}

	_ = c.ServeJSON()
}

// @Title Get
// @Description find group by group id
// @Param	gid		path 	int64	true		"the group id you want to get"
// @Success 200 {object} models.GroupView
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int [get]
func (c *GroupController) GetGroup(gid int64) {
	if group, err := models.GetGroup(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(group.View())
	}

	_ = c.ServeJSON()
}

// @Title GetGroups
// @Description get all groups
// @Param	name		query 	string	false		"filter by group name"
// @Param	page		query 	int	false		"pagination start"
// @Param	pageSize		query 	int	false		"pagination size"
// @Success 200 {object} []models.GroupPreview
// @Failure 500 internal error
// @router / [get]
func (c *GroupController) GetGroups(name *string, page *int, pageSize *int) {
	_page := 1
	_pageSize := 100

	if pageSize != nil && *pageSize < 100 {
		_pageSize = *pageSize
	}

	if page != nil && *page > 0 {
		_page = *page
	}

	if list, err := models.GetGroups(name, _page, _pageSize); err != nil {
		CustomAbort(c, err, 500, "[]")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title Update
// @Description update the group
// @Param	gid		path 	int64	true		"The group id you want to update"
// @Param	body		body 	models.Group	true		"The body"
// @Success 200 {object} models.GroupView
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int [put]
func (c *GroupController) PutGroup(gid int64) {
	if g, _, err := models.Verify[models.Group](c.Ctx.Input.RequestBody, gid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if ug, err := models.UpdateGroup(gid, g); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(ug.View())
	}

	_ = c.ServeJSON()
}

// @Title Delete
// @Description delete the group
// @Param	gid		path 	string	true		"The group id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 invalid gid
// @router /:gid:int [delete]
func (c *GroupController) DeleteGroup(gid int64) {
	// TODO: Requester must be admin, or group admin
	if n, err := models.DeleteGroup(gid); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}

	_ = c.ServeJSON()
}

// @Title Add User
// @Description add user to the group
// @Param	gid		path 	int64	true		"The group id you want to update"
// @Param	uid		path 	int64	true		"The user id you want to add"
// @Param	secret		query 	string	false		"A secret invitation code"
// @Success 200 {int64} user group relation Id
// @Failure 400 invalid input
// @router /:gid:int/user/:uid:int [post]
func (c *GroupController) AddUser(gid int64, uid int64, secret *string) {
	if g, err := models.GetGroup(gid); err != nil {
		CustomAbort(c, nil, 404, "Not Found")
	} else if secret != nil && (g.Token == nil || *g.Token != *secret) {
		CustomAbort(c, nil, 403, "Not Authorized")
	} else if secret == nil {
		// TODO: If Secret == nil requester must be admin or group admin
		CustomAbort(c, nil, 403, "Not Authorized")
	} else if u, err := models.GetUser(uid); err != nil {
		CustomAbort(c, nil, 404, "Not Found")
	} else if rid, err := models.GroupAddUser(g, u, ""); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{"oid": strconv.FormatInt(rid, 10)})
	}

	_ = c.ServeJSON()
}

// @Title List User
// @Description get all users in group
// @Param	gid		path 	string	true		"The group id you want to inspect"
// @Success 200 {object} []models.UserPreview
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int/user/ [get]
func (c *GroupController) GetUsers(gid int64) {
	// TODO: Requester must be part of the group or admin
	if list, err := models.GroupGetUsers(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title Remove User
// @Description add user to the group
// @Param	gid		path 	string	true		"The group id you want to update"
// @Param	uid		path 	string	true		"The user id you want to add"
// @Success 200 {int64} number of deleted rows
// @Failure 400 invalid input
// @Failure 404 gid or uid does not exist
// @router /:gid:int/user/:uid:int [delete]
func (c *GroupController) RemoveUser(gid int64, uid int64) {
	// TODO: Requester must be admin, or group admin, or removed user.
	_, err := models.GroupRemoveUser(gid, uid)
	if err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData("")
	}

	_ = c.ServeJSON()
}

// @Title Add Subgroup
// @Description add subgroup to the group
// @Param	gid		path 	int64	true		"The group id you want to update"
// @Param	sgid		path 	int64	true		"The subgroup id you want to add"
// @Param	secret		query 	string	false		"A secret invitation code"
// @Success 200 {int64} group to subgroup relation Id
// @Failure 400 invalid input
// @router /:gid:int/group/:sgid:int [post]
func (c *GroupController) AddSubgroup(gid int64, sgid int64, secret *string) {
	if g, err := models.GetGroup(gid); err != nil {
		CustomAbort(c, nil, 404, "Not Found")
	} else if secret != nil && (g.Token == nil || *g.Token != *secret) {
		CustomAbort(c, nil, 403, "Not Authorized")
	} else if secret == nil {
		// TODO: If Secret == nil requester must be admin or group admin
		CustomAbort(c, nil, 403, "Not Authorized")
	} else if sg, err := models.GetGroup(sgid); err != nil {
		CustomAbort(c, nil, 404, "Not Found")
	} else if rid, err := models.GroupAddSubgroup(g, sg); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{"oid": strconv.FormatInt(rid, 10)})
	}

	_ = c.ServeJSON()
}

// @Title List Subgroup
// @Description get all subgroups in group
// @Param	gid		path 	string	true		"The group id you want to inspect"
// @Success 200 {object} []models.GroupPreview
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int/sub/ [get]
func (c *GroupController) GetSubgroup(gid int64) {
	// TODO: Requester must be admin, or group admin
	if list, err := models.GroupGetSubgroups(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title List Supgroup
// @Description get all subgroups in group
// @Param	gid		path 	string	true		"The group id you want to inspect"
// @Success 200 {object} []models.GroupPreview
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int/sup/ [get]
func (c *GroupController) GetSupgroup(gid int64) {
	// TODO: Requester must be admin, or group admin
	if list, err := models.GroupGetSupgroups(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title Remove Subgroup
// @Description add sub group to the group
// @Param	gid		path 	string	true		"The group id you want to update"
// @Param	sgid		path 	string	true		"The group id you want to add"
// @Success 200 {int64} number of deleted rows
// @Failure 400 invalid input
// @Failure 404 gid or sgid does not exist
// @router /:gid:int/group/:sgid:int [delete]
func (c *GroupController) RemoveSubgroup(gid int64, sgid int64) {
	// TODO: Requester must be admin, or group admin
	if n, err := models.GroupRemoveSubgroup(gid, sgid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}

	_ = c.ServeJSON()
}
