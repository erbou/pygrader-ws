package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserGroup struct {
	Id      int64  `orm:"pk;auto"`
	User    *User  `orm:"rel(fk)"`
	Group   *Group `orm:"rel(fk)"`
	GrpAcl  string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(UserGroup))
}

func (ug *UserGroup) TableUnique() [][]string {
	return [][]string{
		{"User", "Group"},
	}
}

func (ug *UserGroup) TableName() string {
	return "m2m_user_group"
}

func AddUserGroup(u *User, g *Group, acl string) (*UserGroup, error) {
	o := orm.NewOrm()
	uG := UserGroup{
		User:    u,
		Group:   g,
		GrpAcl:  acl,
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}
	_, err := o.Insert(&uG)
	return &uG, err
}
