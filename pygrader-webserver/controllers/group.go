package controllers

import (
	"pygrader-webserver/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about groups
type GroupController struct {
	beego.Controller
}

func (c *GroupController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title Create
// @Description create group
// @Param	body		body 	models.Group	true		"The group content"
// @Success 200 {object} models.Group.Id
// @Failure 400 invalid input
// @router / [post]
func (c *GroupController) Post() {
	if g, s, err := models.Verify[models.Group](c.Ctx.Input.RequestBody); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if obj, err := models.AddGroup(g); err != nil {
		CustomAbort(c, err, 400, "Failed")
	} else {
		if _, err := models.GroupAddUser(g.Id, s.Issuer.Id, nil); err != nil {

		}
		c.SetData(obj)
	}
	c.ServeJSON()
}

// @Title Get
// @Description find group by group id
// @Param	gid		path 	int64	true		"the group id you want to get"
// @Success 200 {object} models.Group
// @Failure 400 invalid gid
// @Failure 404 gid does not exist
// @router /:gid:int [get]
func (c *GroupController) GetGroup(gid int64) {
	if group, err := models.GetGroup(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(group)
	}
	c.ServeJSON()
}

// @Title GetGroups
// @Description get all groups
// @Param	name		query 	string	false		"filter by group name"
// @Success 200 {object} []models.Group
// @Success 500 internal error
// @router / [get]
func (c *GroupController) GetGroups(name *string) {
	if list, err := models.GetGroups(name); err != nil {
		CustomAbort(c, err, 500, "[]")
	} else {
		c.SetData(list)
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the group
// @Param	gid		path 	int64	true		"The group id you want to update"
// @Param	body		body 	models.Group	true		"The body"
// @Success 200 {object} models.Group
// @Failure 400 invalid input
// @Failure 404 gid does not exist
// @router /:gid:int [put]
func (c *GroupController) PutGroup(gid int64) {
	if g, _, err := models.Verify[models.Group](c.Ctx.Input.RequestBody, gid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if gg, err := models.UpdateGroup(gid, g); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(gg)
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the group
// @Param	gid		path 	string	true		"The group id you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 invalid gid
// @router /:gid:int [delete]
func (c *GroupController) DeleteGroup(gid int64) {
	if n, err := models.DeleteGroup(gid); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}
	c.ServeJSON()
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
	rid, err := models.GroupAddUser(gid, uid, secret)
	if err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{"oid": strconv.FormatInt(rid, 10)})
	}
	c.ServeJSON()
}

// @Title List User
// @Description get all users in group
// @Param	gid		path 	string	true		"The group id you want to inspect"
// @Success 200 {object} []models.User
// @Failure 400 invalid gid
// @Failure 404 gid does not exist
// @router /:gid:int/user/ [get]
func (c *GroupController) GetUsers(gid int64) {
	if users, err := models.GroupGetUsers(gid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(users)
	}
	c.ServeJSON()
}

// @Title Remove User
// @Description add user to the group
// @Param	gid		path 	string	true		"The group id you want to update"
// @Param	uid		path 	string	true		"The user id you want to add"
// @Success 200 {int64} user group relation Id
// @Failure 400 invalid input
// @Failure 404 gid or uid does not exist
// @router /:gid:int/user/:uid:int [delete]
func (c *GroupController) RemoveUser(gid int64, uid int64) {
	_, err := models.GroupRemoveUser(gid, uid)
	if err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData("")
	}
	c.ServeJSON()
}
