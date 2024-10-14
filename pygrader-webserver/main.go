package main

import (
	"regexp"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	_ "pygrader-webserver/routers"
)

const (
	maxIdleConnections = 30
	maxOpenConnections = 30
	asyncLogsTime      = 1e3
)

func init() {
}

func main() {
	// set default database
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		logs.Critical("%v", err)
	}

	force := false
	verbose := false
	beego.BConfig.EnableGzip = true

	if matched, _ := regexp.MatchString("(?i).*dev.*", beego.BConfig.RunMode); matched {
		verbose = true
		force = true
		orm.Debug = true
		logs.EnableFuncCallDepth(true)
		_ = logs.SetLogger(logs.AdapterConsole, `{"level":1}`)
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	} else {
		orm.Debug = false

		logs.Reset()
		_ = logs.SetLogger(logs.AdapterFile, `{"level":4, "filename": "pygrader.log", "perm": "0700" }`)
		logs.Async(asyncLogsTime)
	}

	logs.Info("Mode %v", beego.BConfig.RunMode)

	orm.DefaultTimeLoc = time.UTC

	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		logs.Error(err)
	}

	if err := orm.RegisterDataBase("default", "mysql", sqlConn); err != nil {
		logs.Error(err)
	}

	// orm.RegisterDataBase("default", "sqlite3", "sqlite.db")
	// orm.RegisterDriver("sqlite3", orm.DRSqlite)

	orm.SetMaxIdleConns("default", maxIdleConnections)
	orm.SetMaxOpenConns("default", maxOpenConnections)

	orm.RunCommand()

	if err := orm.RunSyncdb("default", force, verbose); err != nil {
		logs.Error(err)
	}

	// DB- specific (mysql), modify as needed
	// o := orm.NewOrm()
	// _, err = o.Raw("alter table user_groups add constraint cst_unique unique(user_id,group_id)").Exec()

	beego.Run()
}
