package controllers

import (
	"pygrader-webserver/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type QuestionController struct {
	beego.Controller
}

func (c *QuestionController) GetController() *beego.Controller {
	return &c.Controller
}

// Operations about Questions
// Note, question are created via the Module controller.

// @Title Get
// @Description get question by id
// @Param	qid		path 	int64	true		"The id of the question being searched"
// @Success 200 {object} models.QuestionView
// @Failure 400 qid is not int64
// @router /:qid:int [get]
func (c *QuestionController) GetQuestion(qid int64) {
	if obj, err := models.GetQuestion(qid); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(obj.View())
	}
	c.ServeJSON()
}

// @Title Update
// @Description update a question
// @Param	qid		path 	int64	true		"The id of the question being updated"
// @Param	body		body 	models.Question	true		"question object"
// @Success 200 {object} models.QuestionView
// @Failure 400 qid is not int64
// @router /:qid:int [put]
func (c *QuestionController) PutQuestion(qid int64) {
	obj := models.Question{MinScore: -1, MaxTry: -1} // -1 used as nil/unset value in body
	if obj, _, err := models.Verify[models.Question](c.Ctx.Input.RequestBody, &obj, qid); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if uu, err := models.UpdateQuestion(qid, obj); err != nil {
		CustomAbort(c, err, 404, "Not Found")
	} else {
		c.SetData(uu.View())
	}
	c.ServeJSON()
}

// @Title Delete
// @Description delete a question
// @Param	qid		path 	int64	true		"The id of the question being deleted"
// @Success 200
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

// @Title Submit Answer
// @Description answer the question
// @Param	qid		path 	int64	true		"The id of the question being answered"
// @Param	obj		body 	models.GroupAnswerInput	true		"The answer"
// @Success 200 {object} models.GroupAnswer
// @Failure 400 bad request
// @Failure 404 question or group not found
// @router /:qid:int/answer [post]
func (c *QuestionController) PostAnswer(qid int64) {
	if g, s, err := models.Verify[models.GroupAnswerInput](c.Ctx.Input.RequestBody); err != nil {
		CustomAbort(c, err, 403, "Forbidden")
	} else if g.Question != qid {
		CustomAbort(c, err, 400, "Bad Request")
	} else {
		g.Poster = s.Issuer
		if obj, err := g.MapInput(); err != nil {
			CustomAbort(c, err, 400, "Bad Request")
		} else if obj, err := models.AddGroupAnswer(obj); err != nil {
			CustomAbort(c, err, 400, "Bad Request")
		} else {
			c.SetData(obj.View())
		}
	}
	c.ServeJSON()
}
