package models

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"pygrader-webserver/uti"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id       int64
	Username string    `orm:"unique" hash:"n"`
	Email    string    `orm:"unique" hash:"e"`
	Key      string    `orm:"size(4096)" hash:"k"`
	CName    string    `orm:"size(32);unique;index"`
	CEmail   string    `orm:"size(64);unique;index"`
	Kid      string    `orm:"size(40);unique;index"`
	Scope    *string   `orm:"null" hash:"s"`
	Groups   []*Group  `orm:"rel(m2m);rel_through(pygrader-webserver/models.UserGroup)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

type UserPreview struct {
	Id       int64
	Username string
	Email    string
	Kid      string
}

type UserView struct {
	Id       int64
	Username string
	Email    string
	Key      string
	Kid      string
	Scope    *string
	Created  time.Time
	Updated  time.Time
}

func (obj *User) Preview() *UserPreview {
	if obj == nil {
		return nil
	}

	return &UserPreview{
		Id:       obj.Id,
		Username: obj.Username,
		Email:    obj.Email,
		Kid:      obj.Kid,
	}
}

func (obj *User) View() *UserView {
	if obj == nil {
		return nil
	}

	return &UserView{
		Id:       obj.Id,
		Username: obj.Username,
		Email:    obj.Email,
		Key:      obj.Key,
		Kid:      obj.Kid,
		Scope:    obj.Scope,
		Created:  obj.Created,
		Updated:  obj.Updated,
	}
}

func (obj *User) Validate() error {
	if obj == nil {
		return uti.Errorf(ErrInvalidInput, "Invalid input")
	} else if obj.Username = strings.Trim(obj.Username, " \t\n"); obj.Username == "" {
		return uti.Errorf(ErrInvalidInput, "Invalid input '%v'", obj.Username)
	} else if cname, err := uti.CanonizeName(obj.Username); err != nil {
		return err
	} else {
		obj.CName = cname
	}

	if email, err := uti.CanonizeEmail(obj.Email); err != nil {
		return err
	} else {
		obj.CEmail = email
	}

	if kid, key, err := uti.CryptoGetKeyFingerprint(obj.Key, uti.KeyAny, 40); err != nil {
		return err
	} else {
		obj.Kid = kid
		obj.Key = base64.StdEncoding.EncodeToString(key)
	}

	return nil
}

func AddUser(obj *User) (*User, error) {
	// Compute key fingerprint
	if err := obj.Validate(); err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	if _, err := o.Insert(obj); err != nil {
		return nil, err
	} else {
		return obj, err
	}
}

func GetUser(oid int64) (*User, error) {
	u := User{Id: oid}
	o := orm.NewOrm()

	if err := o.Read(&u); err == nil {
		return &u, nil
	} else {
		return nil, err
	}
}

func GetUsers(email *string, username *string, kid *string, page int, pageSize int) (*[]*UserPreview, error) {
	var users []*User

	cond := orm.NewCondition()

	if email != nil {
		if _email, err := uti.CanonizeEmail(*email); err != nil {
			return nil, err
		} else {
			cond = cond.And(`CEmail`, _email)
		}
	}

	if username != nil {
		if _username, err := uti.CanonizeName(*username); err != nil {
			return nil, err
		} else {
			cond = cond.And(`CName`, _username)
		}
	}

	if kid != nil {
		cond = cond.And(`kid`, *kid)
	}

	o := orm.NewOrm()
	qs := o.QueryTable("user")

	if _, err := qs.Limit(pageSize, (page-1)*pageSize).SetCond(cond).All(&users, "Id", "Username", "Email", "Kid"); err != nil {
		return nil, err
	}

	list := make([]*UserPreview, 0, len(users))
	for _, u := range users {
		list = append(list, u.Preview())
	}

	return &list, nil
}

func UpdateUser(oid int64, obj *User) (*User, error) {
	u := User{Id: oid}
	o := orm.NewOrm()
	// race condition?
	if err := o.Read(&u); err == nil {
		if obj.Key != "" {
			u.Key = obj.Key
		}

		if obj.Username != "" {
			u.Username = obj.Username
		}

		if obj.Email != "" {
			u.Email = obj.Email
		}

		if err := u.Validate(); err != nil {
			return nil, err
		}

		if _, err := o.Update(&u); err != nil {
			return nil, err
		}

		return &u, nil
	} else {
		return nil, err
	}
}

func DeleteUser(oid int64) (int64, error) {
	o := orm.NewOrm()

	// Remove keys from cache otherwise it will return the wrong user ID
	// if user is deleted and recreated with the same kid before key expires
	_ = GlCache.ClearAll(context.Background())

	return o.Delete(&User{Id: oid})
}

func UserGetGroups(oid int64) (*[]*Group, error) {
	u := User{Id: oid}
	o := orm.NewOrm()

	if err := o.Read(&u); err != nil {
		return nil, err
	}

	if _, err := o.LoadRelated(&u, "Groups"); err != nil {
		return nil, err
	}

	return &u.Groups, nil
}
