package controllers

import (
	"pygrader-webserver/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Questions
type QuestionController struct {
	beego.Controller
}

func (c *QuestionController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title CreateQuestion
// @Description create question
// @Param	mid		path 	int64	true		"The module id of the question"
// @Param	body		body 	models.Question	true	"question object"
// @Success 200 {int64} models.Question.Id
// @Failure 400 body is invalid
// @router /module/:mid:int [post]
func (c *QuestionController) PostQuestion(mid int64) {
	if iObj, _, err := models.Verify[models.Question](c.Ctx.Input.RequestBody, mid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if mObj, err := models.AddQuestion(mid, iObj); err != nil {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		c.SetData(mObj)
	}
	c.ServeJSON()
}

// @Title Get
// @Description get question by id
// @Param	qid		path 	int64	true		"The question id you want to get"
// @Success 200 {object} models.Question
// @Failure 400 qid is not int64
// @router /:qid:int [get]
func (c *QuestionController) GetQuestion(qid int64) {
	if obj, err := models.GetQuestion(qid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(obj)
	}
	c.ServeJSON()
}

// @Title Update
// @Description update the question
// @Param	qid		path 	int64	true		"The id of the question you want to update"
// @Param	body		body 	models.Question	true		"question object"
// @Success 200 {object} models.Question
// @Failure 400 qid is not int64
// @router /:qid:int [put]
func (c *QuestionController) PutQuestion(qid int64) {
	obj := models.Question{MinScore: -1} // -1 used as nil/unset value in body
	if obj, _, err := models.Verify[models.Question](c.Ctx.Input.RequestBody, &obj, qid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uu, err := models.UpdateQuestion(qid, obj); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(uu)
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete the question
// @Param	qid		path 	int64	true		"The id of the question you want to delete"
// @Success 200 {string} delete success!
// @Failure 400 qid is not int64
// @router /:qid:int [delete]
func (c *QuestionController) DeleteQuestion(qid int64) {
	if n, err := models.DeleteQuestion(qid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(map[string]string{`nrow`: strconv.FormatInt(n, 10)})
	}
	c.ServeJSON()
}
