package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserGroup struct {
	Id      int64  `orm:"pk;auto"`
	User    *User  `orm:"rel(fk);null;on_delete(set_null)"`
	Group   *Group `orm:"rel(fk);null;on_delete(set_null)"`
	GrpACL  string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(UserGroup))
}

func (obj *UserGroup) TableUnique() [][]string {
	return [][]string{
		{"User", "Group"},
	}
}

func (obj *UserGroup) TableName() string {
	return "m2m_user_group"
}

func AddUserGroup(u *User, g *Group, acl string) (*UserGroup, error) {
	o := orm.NewOrm()
	uG := UserGroup{
		User:    u,
		Group:   g,
		GrpACL:  acl,
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}

	if nrow, err := o.InsertOrUpdate(&uG); nrow > 0 {
		return &uG, err
	} else {
		err := o.Raw("SELECT * FROM m2m_user_group WHERE group_id = ? AND user_id = ?", g.Id, u.Id).QueryRow(&uG)

		return &uG, err
	}
}
