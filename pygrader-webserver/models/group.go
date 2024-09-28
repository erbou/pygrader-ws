package models

import (
	//"errors"
	"fmt"
	"pygrader-webserver/uti"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
)

type Group struct {
	Id      int64
	Name    string    `orm:"unique" hash:"n"`
	Scope   *string   `orm:"null" hash:"s"`
	Token   *string   `orm:"null" hash:"t"`
	CName   string    `orm:"unique;size(32)"`
	Users   []*User   `orm:"reverse(many)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Group))
}

func (obj *Group) Validate() error {
	if obj == nil {
		return fmt.Errorf("invalid input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return fmt.Errorf("invalid input")
	} else if cname, err := uti.CanonizeName(obj.Name); err != nil {
		return err
	} else {
		obj.CName = cname
	}
	return nil
}

func AddGroup(obj *Group) (*Group, error) {
	o := orm.NewOrm()
	uuid := uuid.New().String()
	obj.Token = &uuid
	if err := obj.Validate(); err != nil {
		return nil, err
	}
	if n, err := o.Insert(obj); err == nil {
		logs.Debug("Insert group oid:(%v) nrows:(%v)", obj.Id, n)
		return obj, nil
	} else {
		return nil, err
	}
}

func GetGroup(oid int64) (*Group, error) {
	g := Group{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&g); err == nil {
		return &g, nil
	} else {
		return nil, err
	}
}

func GetGroups(name *string) (*[]*Group, error) {
	var list []*Group
	o := orm.NewOrm()
	qs := o.QueryTable("group")

	cond := orm.NewCondition()
	if name != nil {
		cond = cond.And(`name`, name)
	}

	if _, err := qs.SetCond(cond).All(&list); err != nil {
		return nil, err
	}

	// Do not show the secrets
	for i := range list {
		list[i].Token = nil
	}

	return &list, nil
}

func UpdateGroup(oid int64, obj *Group) (*Group, error) {
	if obj == nil {
		return nil, fmt.Errorf("invalid input")
	}
	dbObj := Group{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&dbObj); err != nil {
		return nil, err
	} else {
		if obj.Name != "" {
			dbObj.Name = obj.Name
		}
		if obj.Token != nil {
			if *obj.Token == "" {
				dbObj.Token = nil
			} else if tok, err := uuid.Parse(*obj.Token); err == nil && dbObj.Token != nil && tok.String() == *dbObj.Token {
				uuid := uuid.New().String()
				dbObj.Token = &uuid
			}
		}
		if err := dbObj.Validate(); err != nil {
			return nil, err
		}
		if n, err := o.Update(&dbObj); err != nil {
			return nil, err
		} else {
			logs.Debug("Update user group oid:(%v) nrows:(%v)", dbObj.Id, n)
		}
		return &dbObj, nil
	}
}

func DeleteGroup(oid int64) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&Group{Id: oid})
}

func GroupAddUser(oid int64, uid int64, secret *string) (int64, error) {
	g := Group{Id: oid}
	u := User{Id: uid}
	o := orm.NewOrm()
	if err := o.Read(&g); err != nil {
		return 0, err
	}
	if secret == nil || g.Token == nil || secret != g.Token {
		// Secret token is not being used to join
		// TODO: user must be a group member
		if false {
		} else {
			return 0, fmt.Errorf("Forbidden")
		}
	}
	if err := o.Read(&u); err != nil {
		return 0, err
	}
	uG, err := AddUserGroup(&u, &g, "")
	return uG.Id, err
}

func GroupGetUsers(oid int64) (*[]*User, error) {
	g := Group{Id: oid}
	o := orm.NewOrm()
	if err := o.Read(&g); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(&g, "Users"); err != nil {
		return nil, err
	}
	return &g.Users, nil
}

func GroupRemoveUser(oid int64, uid int64) (int64, error) {
	g := Group{Id: oid}
	u := User{Id: uid}
	o := orm.NewOrm()
	if err := o.Read(&g); err != nil {
		return 0, err
	}
	if err := o.Read(&u); err != nil {
		return 0, err
	}
	m2m := o.QueryM2M(&g, "Users")
	n, err := m2m.Remove(u)
	logs.Debug("Delete user group oid:(%v) nrow:(%v) err:(%v)", oid, n, err)
	return n, err
}
