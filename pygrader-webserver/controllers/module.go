package controllers

import (
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
	"pygrader-webserver/models"
)

// Operations about Modules
type ModuleController struct {
	beego.Controller
}

func (c *ModuleController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title CreateModule
// @Description create module
// @Param	body		body 	models.ModuleInput	true		"module object"
// @Success 200 {object} models.ModuleView
// @Failure 400 bad request
// @router / [post]
func (c *ModuleController) Post() {
	if m, _, err := models.Verify[models.ModuleInput](c.Ctx.Input.RequestBody); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if mm, err := m.MapInput(); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else if uu, err := models.AddModule(mm); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(uu.View())
	}

	_ = c.ServeJSON()
}

// @Title Get
// @Description get all modules
// @Param	name		query 	string	false		"filter by name"
// @Param	page		query 	int	false		"pagination start"
// @Param	pageSize		query 	int	false		"pagination size"
// @Success 200 {object} []models.ModulePreview
// @router / [get]
func (c *ModuleController) GetModules(name *string, page *int, pageSize *int) {
	_page := 1
	_pageSize := 100

	if pageSize != nil && *pageSize < 100 {
		_pageSize = *pageSize
	}

	if page != nil && *page > 0 {
		_page = *page
	}

	if list, err := models.GetModules(name, _page, _pageSize); err != nil {
		CustomAbort(c, err, 500, "System error")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title Get
// @Description get module by mid
// @Param	mid		path 	int64	true		"The module id you want to get"
// @Success 200 {object} models.ModuleView
// @Failure 404 not found
// @router /:mid:int [get]
func (c *ModuleController) GetModule(mid int64) {
	if obj, err := models.GetModule(mid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(obj.View())
	}

	_ = c.ServeJSON()
}

// @Title Update
// @Description update the module
// @Param	mid		path 	int64	true		"The id of the module you want to update"
// @Param	body		body 	models.ModuleInput	true		"module object"
// @Success 200 {object} models.ModuleView
// @Failure 404 not found
// @router /:mid:int [put]
func (c *ModuleController) PutModule(mid int64) {
	if m, _, err := models.Verify[models.ModuleInput](c.Ctx.Input.RequestBody, mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if mm, err := m.MapInput(); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else if uu, err := models.UpdateModule(mid, mm); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(uu)
	}

	_ = c.ServeJSON()
}

// @Title Add Question
// @Description add question
// @Param	mid		path 	int64	true		"The module id"
// @Param	body		body 	models.Question	true	"question object"
// @Success 200 {object} models.QuestionView
// @Failure 400 Bad Request
// @Failure 403 Forbidden
// @router /:mid:int/question [post]
func (c *ModuleController) PostQuestion(mid int64) {
	if iObj, _, err := models.Verify[models.Question](c.Ctx.Input.RequestBody, mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if mObj, err := models.AddQuestion(mid, iObj); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(mObj.View())
	}

	_ = c.ServeJSON()
}

// @Title Delete
// @Description delete the module
// @Param	mid		path 	int64	true		"The mid you want to delete"
// @Success 200 {int} number of rows deleted
// @Failure 403 forbidden
// @router /:mid:int [delete]
func (c *ModuleController) DeleteModule(mid int64) {
	if n, err := models.DeleteModule(mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}

	_ = c.ServeJSON()
}

// @Title List Questions
// @Description get all questions in module
// @Param	mid		path 	int64	true		"The module id you want to inspect"
// @Success 200 {object} []models.Question
// @Failure 404 not found
// @router /:mid:int/question/ [get]
func (c *ModuleController) GetAllQuestions(mid int64) {
	if list, err := models.ModuleGetAllQuestions(mid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}

// @Title List Answers
// @Description get all answers in module
// @Param	mid		path 	int64	true		"The module id you want to inspect"
// @Param	question		query 	int64	false		"A question id"
// @Param	group		query 	int64	false		"A group id"
// @Param	poster		query 	int64	false		"A user id"
// @Param	page		query 	int	false		"pagination start"
// @Param	pageSize		query 	int	false		"pagination size"
// @Success 200 {object} []models.GroupAnswer
// @Failure 404 not found
// @router /:mid:int/answer/ [get]
func (c *ModuleController) GetGroupAnswers(mid int64, question *int64, group *int64, poster *int64, page *int, pageSize *int) {
	_page := 1
	_pageSize := 100

	if pageSize != nil && *pageSize < 100 {
		_pageSize = *pageSize
	}

	if page != nil && *page > 0 {
		_page = *page
	}

	if list, err := models.GetGroupAnswers(mid, question, group, poster, _page, _pageSize); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(list)
	}

	_ = c.ServeJSON()
}
