package models

import (
	//"errors"
	"pygrader-webserver/uti"
	"strings"
	"time"

	orm "github.com/beego/beego/v2/client/orm"
)

type Question struct {
	Id       int64
	Name     string     `orm:"size(32"             hash:"n"`
	Before   *time.Time `orm:"type(datetime);null" hash:"b"`
	Grader   string     `orm:"size(32)"            hash:"g"`
	MaxScore int        `hash:"h"`
	MinScore int        `hash:"m"`
	CName    string     `orm:"size(32)"`
	Module   *Module    `orm:"rel(fk);on_delete(do_nothing)"`
	Created  time.Time  `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time  `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Question))
}

func (ug *Question) TableUnique() [][]string {
	return [][]string{
		{"CName", "Module"},
	}
}

func (obj *Question) Validate() error {
	if obj == nil {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if obj.MinScore < 0 || obj.MaxScore <= obj.MinScore {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid min/max score (%v,%v)", obj.MinScore, obj.MaxScore)
	} else if cname, err := uti.CanonizeName(obj.Name); err != nil {
		return err
	} else {
		obj.CName = cname
	}
	return nil
}

func AddQuestion(oid int64, obj *Question) (*Question, error) {
	o := orm.NewOrm()
	if err := obj.Validate(); err != nil {
		return nil, err
	} else if m, err := GetModule(oid); err != nil {
		return nil, err
	} else {
		obj.Module = m
	}

	if _, err := o.Insert(obj); err != nil {
		return nil, err
	} else {
		return obj, nil
	}
}

func GetQuestion(oid int64) (*Question, error) {
	obj := Question{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&obj); err == nil {
		return &obj, nil
	} else {
		return nil, err
	}
}

func GetAllQuestions() (*[]*Question, error) {
	var list []*Question
	o := orm.NewOrm()
	qs := o.QueryTable("question")
	if _, err := qs.All(&list); err == nil {
		return &list, nil
	} else {
		return nil, err
	}
}

func UpdateQuestion(oid int64, obj *Question) (*Question, error) {
	dbObj := Question{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&dbObj); err == nil {
		if obj.Name != "" {
			dbObj.Name = obj.Name
		}
		if obj.Grader != "" {
			dbObj.Grader = obj.Grader
		}
		if obj.MaxScore != 0 {
			dbObj.MaxScore = obj.MaxScore
		}
		if obj.MinScore >= 0 {
			dbObj.MinScore = obj.MinScore
		}
		if obj.Before != nil {
			dbObj.Before = obj.Before
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

func DeleteQuestion(oid int64) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&Question{Id: oid})
}
