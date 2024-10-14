package controllers

import (
	"errors"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"pygrader-webserver/models"
	"pygrader-webserver/uti"
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
// @Param	body		body 	models.UserInput	true		"user object"
// @Success 200 {object} models.UserView
// @Failure 400 bad request
// @Failure 403 forbidden
// @router / [post]
func (c *UserController) Post() {
	u, s, err := models.Verify[models.User](c.Ctx.Input.RequestBody)
	if err != nil {
		// TODO: Improve bootstraping if kid does not exist yet.
		//       Current approach is to allow it if the new user's email u.Email and the key fingerprint u.Kid matches the
		//       certificate's CN and signature's key s.Keyid. This requires doing a user Verify pre-validation to obtain
		//       the Kid.
		var _err *uti.Error
		if !errors.As(err, &_err) {
			CustomAbort(c, err, 500, "System error")
		} else if _err.Code == models.ErrKidUnknown &&
			u.Validate() == nil &&
			s.Keyid == u.Kid &&
			c.Ctx.Request.Header.Get("X-Client-Cert-Cn") == u.Email {
			logs.Info("Create user %v %v %v", c.Ctx.Request.Header.Get("X-Client-Cert-Cn"), u.CEmail, u.Kid)
		} else {
			CustomAbort(c, err, 403, "Forbidden")
		}
	}

	if u, err := models.AddUser(u); err != nil {
		CustomAbort(c, err, 400, "Bad request")
	} else {
		c.SetData(u.View())
	}

	_ = c.ServeJSON()
}

// @Title GetUsers
// @Description get all Users
// @Param	email		query 	string	false		"filter by email"
// @Param	username		query 	string	false		"filter by username"
// @Param	kid		query 	string	false		"filter by key fingerprint"
// @Param	page		query 	int	false		"pagination start"
// @Param	pageSize		query 	int	false		"pagination size"
// @Success 200 {object} []models.UserPreview
// @Failure 400 bad request
// @Failure 403 forbidden
// @Failure 500 system error
// @router / [get]
func (c *UserController) GetUsers(email *string, username *string, kid *string, page *int, pageSize *int) {
	_page := 1
	_pageSize := 100

	if pageSize != nil && *pageSize < 100 {
		_pageSize = *pageSize
	}

	if page != nil && *page > 0 {
		_page = *page
	}

	if list, err := models.GetUsers(email, username, kid, _page, _pageSize); err != nil {
		CustomAbort(c, err, 500, "System error")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title Get User
// @Description get user by uid
// @Param	uid		path 	int64	true		"The user id you want to get"
// @Success 200 {object} models.UserView
// @Failure 400 bad request
// @Failure 404 not found
// @router /:uid:int [get]
func (c *UserController) GetUser(uid int64) {
	if u, err := models.GetUser(uid); err == nil {
		c.SetData(u.View())
	} else {
		CustomAbort(c, err, 404, "Not Found")
	}

	_ = c.ServeJSON()
}

// @Title Update User
// @Description update the user
// @Param	uid		path 	int64	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"user object"
// @Success 200 {object} models.User
// @Failure 400 bad request
// @Failure 403 forbidden
// @router /:uid:int [put]
func (c *UserController) PutUser(uid int64) {
	if u, s, err := models.Verify[models.User](c.Ctx.Input.RequestBody, uid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uid != s.Issuer.Id {
		CustomAbort(c, err, 403, "Forbidden")
	} else if u, err := models.UpdateUser(uid, u); err != nil {
		CustomAbort(c, err, 400, "Bad request")
	} else {
		c.SetData(u.View())
	}

	_ = c.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	int64	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 bad request
// @router /:uid:int [delete]
func (c *UserController) DeleteUser(uid int64) {
	// TODO: Only admin can delete a user
	if n, err := models.DeleteUser(uid); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}

	_ = c.ServeJSON()
}

// @Title List Groups
// @Description get all groups that user belongs to
// @Param	uid		path 	int64	true		"The user id you want to inspect"
// @Success 200 {object} []models.Group
// @Failure 400 bad request
// @Failure 404 uid does not exist
// @router /:uid:int/group/ [get]
func (c *UserController) GetGroups(uid int64) {
	if list, err := models.UserGetGroups(uid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}
