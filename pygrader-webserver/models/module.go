package models

import (
	//"errors"
	"pygrader-webserver/uti"
	"strings"
	"time"

	orm "github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(Module))
}

type Module struct {
	Id        int64
	Name      string      `orm:"unique"`
	CName     string      `orm:"unique;size(32)"`
	Before   *time.Time   `orm:"type(datetime);null"`
	Reveal   *time.Time   `orm:"type(datetime);null"`
	Audience  *Group      `orm:"rel(fk);null;on_delete(set_null)"`
	Questions []*Question `orm:"reverse(many)"`
	Created   time.Time   `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time   `orm:"auto_now;type(datetime)"`
}

type ModuleInput struct {
	Name      string      `hash:"n"`
	Audience  int64       `hash:"a"`
	Before   *time.Time   `hash:"b"`
	Reveal   *time.Time   `hash:"r"`
}

type ModulePreview struct {
	Id        int64
	Name      string
	Before   *time.Time
}

type ModuleView struct {
	Id        int64
	Name      string
	Audience *GroupPreview
	Before   *time.Time
	Reveal   *time.Time
	Created   time.Time
	Updated   time.Time
}

func (obj *ModuleInput) MapInput() (*Module, error) {
	if audience, err := GetGroup(obj.Audience); err != nil {
		return nil, err
	} else {
		return &Module{Name: obj.Name, Before: obj.Before, Reveal: obj.Reveal, Audience: audience}, nil
	}
}

func (obj *Module) Preview() *ModulePreview {
	if obj == nil {
		return nil
	}
	return &ModulePreview{Id: obj.Id, Name: obj.Name, Before: obj.Before}
}

func (obj *Module) View() *ModuleView {
	if obj == nil {
		return nil
	}
	return &ModuleView{
		Id: obj.Id,
		Name: obj.Name,
		Audience: obj.Audience.Preview(),
		Before: obj.Before,
		Reveal: obj.Reveal,
		Created: obj.Created,
		Updated: obj.Updated,
	}
}

func (obj *Module) Validate() error {
	if obj == nil {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if cname, err := uti.CanonizeName(obj.Name); err != nil {
		return err
	} else if obj.Reveal == nil || obj.Before == nil || obj.Reveal.Before(*obj.Before) || time.Now().After(*obj.Before) {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
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

func GetModules(name *string, page int, pageSize int) (*[]*ModulePreview, error) {
	var list []*Module
	o := orm.NewOrm()
	qs := o.QueryTable("module")

	if _, err := qs.Limit(pageSize, (page-1)*pageSize).All(&list, "Id", "Name"); err != nil {
		return nil, err
	}

	modules := make([]*ModulePreview, 0, len(list))
	for _, m := range list {
		modules = append(modules, m.Preview())
	}

	return &modules, nil
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
		if obj.Audience != nil {
			dbObj.Audience = obj.Audience
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
