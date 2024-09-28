package controllers

import (
	"pygrader-webserver/models"
	"pygrader-webserver/uti"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

func (c *UserController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"user object"
// @Success 200 {int64} models.User.Id
// @Failure 400 body is invalid
// @router / [post]
func (c *UserController) Post() {
	u, s, err := models.Verify[models.User](c.Ctx.Input.RequestBody)
	if err != nil {
		// TODO: address bootstraping issue if kid does not exist yet.
		//       current approach is to allow it if the new user's email and the key fingerprint matches the
		//       certificate's CN and signature's key. This requires doing a user pre-validation to obain
		//       the Kid. None of this is very clean.
		if _err, ok := err.(uti.Error); !ok {
			CustomAbort(c, err, 500, "System error")
		} else if _err.Code == models.ERR_KID_UNKOWN &&
			u.Validate() == nil &&
			s.Keyid == u.Kid &&
			c.Ctx.Request.Header.Get("X-Client-Cert-CN") == u.Email {
			logs.Warning("Create user %v %v %v", c.Ctx.Request.Header.Get("X-Client-Cert-CN"), u.CEmail, u.Kid)
		} else {
			CustomAbort(c, err, 403, "Forbidden")
		}
	}

	if u, err := models.AddUser(u); err != nil {
		CustomAbort(c, err, 400, "Failed")
	} else {
		c.SetData(u)
	}
	c.ServeJSON()
}

// @Title GetUsers
// @Description get all Users
// @Param	email		query 	string	false		"filter by email"
// @Param	username		query 	string	false		"filter by username"
// @Param	kid		query 	string	false		"filter by key fingerprint"
// @Success 200 {object} models.User
// @router / [get]
func (c *UserController) GetUsers(email *string, username *string, kid *string) {
	if list, err := models.GetUsers(email, username, kid); err != nil {
		CustomAbort(c, err, 500, "Error")
	} else {
		c.SetData(list)
	}
	c.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	int64	true		"The user id you want to get"
// @Success 200 {object} models.User
// @Failure 400 uid is not int64
// @router /:uid:int [get]
func (c *UserController) GetUser(uid int64) {
	if u, err := models.GetUser(uid); err == nil {
		c.SetData(u)
	} else {
		CustomAbort(c, err, 404, "Not Found")
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	int64	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"user object"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int64
// @router /:uid:int [put]
func (c *UserController) PutUser(uid int64) {
	if u, s, err := models.Verify[models.User](c.Ctx.Input.RequestBody, uid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uid != s.Issuer.Id {
		CustomAbort(c, err, 403, "Forbidden")
	} else if u, err := models.UpdateUser(uid, u); err != nil {
		CustomAbort(c, err, 400, "Failed")
	} else {
		c.SetData(u)
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	int64	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 uid is not int64
// @router /:uid:int [delete]
func (c *UserController) DeleteUser(uid int64) {
	if n, err := models.DeleteUser(uid); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}
	c.ServeJSON()
}

// @Title List Groups
// @Description get all groups that user belongs to
// @Param	uid		path 	int64	true		"The user id you want to inspect"
// @Success 200 {object} []models.Group
// @Failure 400 invalid uid
// @Failure 404 uid does not exist
// @router /:uid:int/group/ [get]
func (c *UserController) GetGroups(uid int64) {
	if list, err := models.UserGetGroups(uid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}
	c.ServeJSON()
}
