package controllers

import (
	"pygrader-webserver/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Modules
type ModuleController struct {
	beego.Controller
}

func (c *ModuleController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title CreateModule
// @Description create modullist
// @Param	body		body 	models.Module	true		"module object"
// @Success 200 {int64} models.Module.Id
// @Failure 400 body is invalid
// @router / [post]
func (c *ModuleController) Post() {
	if m, _, err := models.Verify[models.Module](c.Ctx.Input.RequestBody); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uu, err := models.AddModule(m); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(uu)
	}
	c.ServeJSON()
}

// @Title Get
// @Description get all Module
// @Success 200 {object} models.Module
// @router / [get]
func (c *ModuleController) Get() {
	if list, err := models.GetAllModules(); err != nil {
		CustomAbort(c, err, 500, "Internal Error")
	} else {
		c.SetData(list)
	}
	c.ServeJSON()
}

// @Title Get
// @Description get module by mid
// @Param	mid		path 	int64	true		"The module id you want to get"
// @Success 200 {object} models.Module
// @Failure 400 mid is not int64
// @router /:mid:int [get]
func (c *ModuleController) GetModule(mid int64) {
	if obj, err := models.GetModule(mid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(obj)
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the module
// @Param	mid		path 	int64	true		"The mid you want to update"
// @Param	body		body 	models.Module	true		"module object"
// @Success 200 {object} models.Module
// @Failure 400 :mid is not int64
// @router /:mid:int [put]
func (c *ModuleController) PutModule(mid int64) {
	if m, _, err := models.Verify[models.Module](c.Ctx.Input.RequestBody, mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uu, err := models.UpdateModule(mid, m); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(uu)
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the module
// @Param	mid		path 	int64	true		"The mid you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 mid is not int64
// @router /:mid:int [delete]
func (c *ModuleController) DeleteModule(mid int64) {
	if n, err := models.DeleteModule(mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}
	c.ServeJSON()
}

// @Title List Questions
// @Description get all questions of module
// @Param	mid		path 	int64	true		"The module id you want to inspect"
// @Success 200 {object} []models.Question
// @Failure 400 :mid is not int64
// @Failure 404 :mid does not exist
// @router /:mid:int/question/ [get]
func (c *ModuleController) GetAllQuestions(mid int64) {
	if list, err := models.ModuleGetAllQuestions(mid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}
	c.ServeJSON()
}
