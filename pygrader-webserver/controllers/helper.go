package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"

	beego "github.com/beego/beego/v2/server/web"
)

type ControllerInterface interface {
	GetController() *beego.Controller
}

func CustomAbort(c ControllerInterface, err error, code int, msg string) {
	if err == nil {
		c.GetController().CustomAbort(code, msg)
	} else {
		v := reflect.ValueOf(err)
		logs.Warning("ERROR -- %v, %v", v.String(), err.Error())
		if err1, ok := err.(*mysql.MySQLError); ok {
			c.GetController().CustomAbort(400, fmt.Sprintf(`{ "error": "%v", "code": %v }`, err1.Error(), err1.Number))
		} else if err2, ok := err.(*sqlite3.Error); ok {
			c.GetController().CustomAbort(400, fmt.Sprintf(`{ "error": "%v", "code": %v }`, err2.Error(), err2.Code))
		} else if errors.Is(err, sql.ErrNoRows) {
			c.GetController().CustomAbort(404, `{ "error": "Not found" }`)
		} else {
			c.GetController().CustomAbort(code, fmt.Sprintf(`{ "error": "%v" }`, msg))
		}
	}
}
