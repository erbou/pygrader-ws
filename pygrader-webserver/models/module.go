package models

import (
	//"errors"
	"pygrader-webserver/uti"
	"strings"
	"time"

	orm "github.com/beego/beego/v2/client/orm"
)

type Module struct {
	Id        int64
	Name      string      `orm:"unique" hash:"n"`
	CName     string      `orm:"unique;size(32)"`
	Questions []*Question `orm:"reverse(many)"`
	Created   time.Time   `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time   `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Module))
}

func (obj *Module) Validate() error {
	if obj == nil {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if cname, err := uti.CanonizeName(obj.Name); err != nil {
		return err
	} else {
		obj.CName = cname
	}
	return nil
}

func AddModule(obj *Module) (*Module, error) {
	o := orm.NewOrm()
	if err := obj.Validate(); err != nil {
		return nil, err
	}
	if _, err := o.Insert(obj); err == nil {
		return obj, nil
	} else {
		return nil, err
	}
}

func GetModule(oid int64) (*Module, error) {
	obj := Module{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&obj); err == nil {
		return &obj, nil
	} else {
		return nil, err
	}
}

func GetAllModules() (*[]*Module, error) {
	var list []*Module
	o := orm.NewOrm()
	qs := o.QueryTable("module")
	if _, err := qs.All(&list); err == nil {
		return &list, nil
	} else {
		return nil, err
	}
}

func UpdateModule(oid int64, obj *Module) (*Module, error) {
	if obj == nil {
		return nil, uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	}
	dbObj := Module{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&dbObj); err == nil {
		if obj.Name != "" {
			dbObj.Name = obj.Name
		}
		if err := dbObj.Validate(); err != nil {
			return nil, err
		}
		if _, err := o.Update(&dbObj); err == nil {
			return &dbObj, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ModuleGetAllQuestions(oid int64) (*[]*Question, error) {
	obj := Module{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&obj); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(&obj, "Questions"); err != nil {
		return nil, err
	}
	return &obj.Questions, nil
}

func DeleteModule(oid int64) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&Module{Id: oid})
}
