package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"pygrader-webserver/uti"
)

type ControllerInterface interface {
	GetController() *beego.Controller
}

func CustomAbort(c ControllerInterface, err error, code int, msg string) {
	if err == nil {
		c.GetController().CustomAbort(code, fmt.Sprintf(`{ "status", "%v", "code": %v }`, msg, code))
	} else {
		v := reflect.ValueOf(err)
		logs.Warning("%v, %v", v.String(), err.Error())

		var (
			err0 *uti.Error
			err1 *mysql.MySQLError
			err2 *sqlite3.Error
		)

		if errors.As(err, &err0) {
			c.GetController().CustomAbort(400, fmt.Sprintf(`{ "status": "%v", "code": %v }`, err0.Error(), int(err0.Code)))
		} else if errors.Is(err, sql.ErrNoRows) {
			c.GetController().CustomAbort(404, `{ "status": "Not Found", "code": 404 }`)
		} else if errors.As(err, &err1) {
			c.GetController().CustomAbort(400, fmt.Sprintf(`{ "status": "%v", "code": %v }`, err1.Error(), err1.Number))
		} else if errors.As(err, &err2) {
			c.GetController().CustomAbort(400, fmt.Sprintf(`{ "status": "%v", "code": %v }`, err2.Error(), err2.Code))
		} else {
			c.GetController().CustomAbort(code, fmt.Sprintf(`{ "status": "%v", "code": %v }`, msg, code))
		}
	}
}
