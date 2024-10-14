package models

import (
	"strings"
	"time"

	orm "github.com/beego/beego/v2/client/orm"
	"pygrader-webserver/uti"
)

func init() {
	orm.RegisterModel(new(Question))
}

type Question struct {
	Id       int64
	Name     string     `orm:"size(32" hash:"n"`
	Before   *time.Time `orm:"type(datetime);null" hash:"b"`
	Reveal   *time.Time `orm:"type(datetime);null" hash:"r"`
	Grader   string     `orm:"size(32)" hash:"g"`
	MaxScore int        `hash:"h"`
	MinScore int        `hash:"m"`
	MaxTry   int        `hash:"t"`
	CName    string     `orm:"size(32)"`
	Module   *Module    `orm:"rel(fk);null;on_delete(set_null)"`
	Created  time.Time  `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time  `orm:"auto_now;type(datetime)"`
}

type QuestionPreview struct {
	Id     int64
	Name   string
	Before *time.Time
}

type QuestionView struct {
	Id       int64
	Name     string
	Before   *time.Time
	Reveal   *time.Time
	Grader   string
	MaxScore int
	MinScore int
	MaxTry   int
	Module   *ModulePreview
	Created  time.Time
	Updated  time.Time
}

func (obj *Question) TableUnique() [][]string {
	return [][]string{
		{"CName", "Module"},
	}
}

func (obj *Question) Validate() error {
	if obj == nil {
		return uti.Errorf(ErrInvalidInput, "Invalid input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return uti.Errorf(ErrInvalidInput, "Invalid input")
	} else if obj.MinScore < 0 || obj.MaxScore <= obj.MinScore {
		return uti.Errorf(ErrInvalidInput, "Invalid min/max score (%v, %v)", obj.MinScore, obj.MaxScore)
	} else if obj.MaxTry < 0 {
		return uti.Errorf(ErrInvalidInput, "Invalid max attempts (%v)", obj.MaxTry)
	} else if obj.Reveal == nil || obj.Before == nil || obj.Reveal.Before(*obj.Before) || time.Now().After(*obj.Before) {
		return uti.Errorf(ErrInvalidInput, "Invalid before and reveal date (%v, %v)", obj.Before, obj.Reveal)
	} else if obj.Module == nil || obj.Before.After(*obj.Module.Before) || obj.Reveal.After(*obj.Module.Reveal) {
		return uti.Errorf(ErrInvalidInput, "Invalid before and reveal are past module date")
	} else if cname, err := uti.CanonizeName(obj.Name); err != nil {
		return err
	} else {
		obj.CName = cname
	}

	return nil
}

func (obj *Question) Preview() *QuestionPreview {
	return &QuestionPreview{Id: obj.Id, Name: obj.Name, Before: obj.Before}
}

func (obj *Question) View() *QuestionView {
	if obj == nil {
		return nil
	}

	return &QuestionView{
		Id:       obj.Id,
		Name:     obj.Name,
		Grader:   obj.Grader,
		MinScore: obj.MinScore,
		MaxScore: obj.MaxScore,
		MaxTry:   obj.MaxTry,
		Before:   obj.Before,
		Reveal:   obj.Reveal,
		Module:   obj.Module.Preview(),
	}
}

func AddQuestion(oid int64, obj *Question) (*Question, error) {
	o := orm.NewOrm()

	if m, err := GetModule(oid); err != nil {
		return nil, err
	} else {
		obj.Module = m
	}

	if err := obj.Validate(); err != nil {
		return nil, err
	} else if n, err := o.Insert(obj); err != nil || n < 1 {
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

/*
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
*/

func UpdateQuestion(oid int64, obj *Question) (*Question, error) {
	o := orm.NewOrm()
	dbObj := Question{Id: oid}

	if err := o.Read(&dbObj); err != nil {
		return nil, err
	}

	if obj.Name != "" {
		dbObj.Name = obj.Name
	}

	if obj.Grader != "" {
		dbObj.Grader = obj.Grader
	}

	if obj.MaxScore > 0 {
		dbObj.MaxScore = obj.MaxScore
	}

	if obj.MinScore >= 0 {
		dbObj.MinScore = obj.MinScore
	}

	if obj.MaxTry >= 0 {
		dbObj.MaxTry = obj.MaxTry
	}

	if obj.Before != nil {
		dbObj.Before = obj.Before
	}

	if obj.Reveal != nil {
		dbObj.Reveal = obj.Reveal
	}

	if err := dbObj.Validate(); err != nil {
		return nil, err
	} else if n, err := o.Update(&dbObj); err != nil || n < 1 {
		return nil, err
	} else {
		return &dbObj, nil
	}
}

func DeleteQuestion(oid int64) (int64, error) {
	o := orm.NewOrm()

	return o.Delete(&Question{Id: oid})
}
