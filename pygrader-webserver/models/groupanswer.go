package models

import (
	"time"
	"pygrader-webserver/uti"

	orm "github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(GroupAnswer))
}

type GroupAnswer struct {
	Id       int64      `orm:"auto"`
	Group    *Group     `orm:"rel(fk);null;on_delete(set_null)"`
	Poster   *User      `orm:"rel(fk);null;on_delete(set_null)"`
	Answer   *Answer    `orm:"rel(fk);null;on_delete(set_null)"`
	NumTry   int
	Created  time.Time  `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time  `orm:"auto_now;type(datetime)"`
}

type GroupAnswerInput struct {
	Question int64      `hash:"q"`
	Group    int64      `hash:"g"`
	Data     string     `hash:"d"`
	Poster   *User
}

type GroupAnswerView struct {
	Group    *GroupPreview
	Poster   *UserPreview
	Answer   *AnswerView
	NumTry   int
	Created  time.Time
	Updated  time.Time
}

func (ug *GroupAnswer) TableName() string {
	return "m2m_group_answer"
}

func (ug *GroupAnswer) TableUnique() [][]string {
	return [][]string{
		{"Poster", "Answer"},
	}
}

func (ug *GroupAnswer) TableIndex() [][]string {
	return [][]string{
		{"Group", "Answer"},
	}
}

func (ug *GroupAnswer) View() *GroupAnswerView {
	return &GroupAnswerView{
		Group: ug.Group.Preview(),
		Poster: ug.Poster.Preview(),
		Answer: ug.Answer.View(),
		NumTry: ug.NumTry,
		Created: ug.Created,
		Updated: ug.Updated,
	}
}

func (obj *GroupAnswerInput) MapInput() (*GroupAnswer, error) {
	if group, err := GetGroup(obj.Group); err != nil {
		return nil, err
	} else if question, err := GetQuestion(obj.Question); err != nil {
		return nil, err
	} else if ans, err := AddAnswer(&Answer{Question: question, Data: obj.Data}); err != nil {
		return nil, err
	} else {
		return &GroupAnswer{Group: group, Poster: obj.Poster, Answer: ans}, nil
	}
}

func (obj *GroupAnswer) Validate() error {
	if obj.Group == nil || obj.Poster == nil || obj.Answer == nil {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	}
	if obj.Answer.Question == nil || obj.Answer.Question.Module == nil || obj.Answer.Question.Module.Audience == nil {
		o := orm.NewOrm()
		if obj.Answer.Question == nil {
			if err := o.Read(obj.Answer); err != nil || obj.Answer.Question == nil {
				return err
			}
		}
		if obj.Answer.Question.Module == nil {
			if err := o.Read(obj.Answer.Question); err != nil || obj.Answer.Question.Module == nil {
				return err
			}
		}
		if obj.Answer.Question.Module.Audience == nil {
			if err := o.Read(obj.Answer.Question.Module); err != nil || obj.Answer.Question.Module.Audience == nil {
				return err
			}
		}
	}
	if time.Now().After(*obj.Answer.Question.Before) {
		return uti.Errorf(ERR_DEADLINE, "Past Deadline")
	}
	return nil
}

func AddGroupAnswer(g *GroupAnswer) (*GroupAnswer, error) {
	o := orm.NewOrm()
	var num_try int64 = 0
	if err := g.Validate(); err != nil {
		return nil, err
	} else if g.Answer.Question.Module.Audience.Id == g.Group.Id {
		// This is a single participant answer
		// TODO: participant must be a member of audience
		if err := o.QueryTable("m2m_group_answer").Filter("Poster__Id", g.Poster.Id).Filter("Answer__Id", g.Answer.Id).One(g); err == nil {
			return g, nil
		}
	} else {
		// This is a group answer
		// TODO: group must be a subgroup of audience
		if err := o.QueryTable("m2m_group_answer").Filter("Group__Id", g.Group.Id).Filter("Answer__Id", g.Answer.Id).One(g); err == nil {
			return g, nil
		}
		if n, err := o.QueryTable("m2m_group_answer").Filter("Group__Id", g.Poster.Id).Filter("Answer__Question__Id", g.Answer.Question.Id).Count(); err != nil {
			return nil, uti.Errorf(uti.ERR_SYSTEM_ERROR, "Try Again")
		} else if n > num_try {
			num_try = n 
		}
	}
	if n, err := o.QueryTable("m2m_group_answer").Filter("Poster__Id", g.Poster.Id).Filter("Answer__Question__Id", g.Answer.Question.Id).Count(); err != nil {
		return nil, uti.Errorf(uti.ERR_SYSTEM_ERROR, "Try Again")
	} else if n > num_try {
		num_try = n 
	}
	g.NumTry = int(num_try + 1)
	if g.Answer.Question.MaxTry > 0 && g.Answer.Question.MaxTry < g.NumTry {
		return nil, uti.Errorf(ERR_MAX_TRY, "Too many attempts")
	}
	if n, err := o.Insert(g); err != nil || n < 1 {
		return nil, uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	}
	return g, nil
}

func GetGroupAnswer(oid int64) (*GroupAnswer, error) {
	obj := GroupAnswer{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&obj); err == nil {
		return &obj, nil
	} else {
		return nil, err
	}
}

func GetGroupAnswers(module int64, question *int64, group *int64, poster *int64, page int, pageSize int) (*[]*GroupAnswer, error) {
	var list []*GroupAnswer
	cond := orm.NewCondition()
	if question != nil {
		cond = cond.And(`Answer__Question__Id`, *question)
	} else {
		cond = cond.And(`Answer__Question__Module_Id`, module)
	}
	if group != nil {
		cond = cond.And(`Group__Id`, *group)
	}
	if poster != nil {
		cond = cond.And(`Poster__Id`, *poster)
	}
	o := orm.NewOrm()
	qs := o.QueryTable("m2m_group_answer")
	//if _, err := qs.Limit(pageSize, (page-1)*pageSize).Filter("Group__Id", group).Filter("Answer__Question__Id", question).All(&list); err == nil {
	if _, err := qs.Limit(pageSize, (page-1)*pageSize).SetCond(cond).All(&list); err == nil {
		return &list, nil
	} else {
		return nil, err
	}
}
