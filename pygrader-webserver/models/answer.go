package models

import (
	"time"
	"pygrader-webserver/uti"

	orm "github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(Answer))
}

type Answer struct {
	Id       int64      `orm:"auto"`
	Data     string     `orm:"type(text)" hash:"a"`
	Result   *string    `orm:"type(text);null" hash:"r"`
	Digest   string     `orm:"size(64)"`
	Score    int        `orm:"default:0"  hash:"s"`
	Question *Question  `orm:"rel(fk);null;on_delete(set_null)" hash:"q"`
	Created  time.Time  `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time  `orm:"auto_now;type(datetime)"`
}

type AnswerView struct {
	Id       int64
	Data     string
	Result   *string
	Digest   string
	Score    int
	Question *QuestionView
	Created  time.Time
	Updated  time.Time
}

func (obj *Answer) TableUnique() [][]string {
	return [][]string{
		{"Digest", "Question"},
	}
}

func (obj *Answer) View() *AnswerView {
	if obj == nil {
		return nil
	}
	return &AnswerView{
		Id: obj.Id,
		Data: obj.Data,
		Digest: obj.Digest,
		Result: obj.Result,
		Score: obj.Score,
		Question: obj.Question.View(),
		Created: obj.Created,
		Updated: obj.Updated,
	}
}

func (obj *Answer) Validate() error {
	if obj.Question == nil || time.Now().After(*obj.Question.Before) {
		return uti.Errorf(ERR_DEADLINE, "Past Deadline")
	}
	return nil
}

func AddAnswer(obj *Answer) (*Answer, error) {
	obj.Digest = uti.HexDigest(obj.Data)
	o := orm.NewOrm()
	if err := obj.Validate(); err != nil {
		return nil, uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if err := o.QueryTable("answer").Filter("Digest", obj.Digest).Filter("Question__Id", obj.Question.Id).One(obj); err != nil {
		if n, err := o.Insert(obj); err != nil || n < 1 {
			return nil, err
		}
	}
	return obj, nil
}

