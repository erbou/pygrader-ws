package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
)

// Operations about groups
type FileController struct {
	beego.Controller
}

func init() {
}

func (c *FileController) GetController() *beego.Controller {
	return &c.Controller
}

// @Title Upload file
// @Description upload a file
// @Success 200
// @Failure 400 invalid input
// @router / [post]
func (c *FileController) Post() {
	f, h, err := c.GetFile("upload")
	uuid := uuid.New().String()
	if err != nil {
		logs.Error("Getfile err ", err)
		CustomAbort(c, err, 400, "Bad Request")
	}
	logs.Debug("Upload ", h)
	defer f.Close()
	c.SaveToFile("upload", "static/upload/" + uuid) 
}
