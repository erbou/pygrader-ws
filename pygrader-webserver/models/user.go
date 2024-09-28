package models

import (
	"encoding/base64"
	"pygrader-webserver/uti"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id       int64
	Username string    `orm:"unique" hash:"n"`
	Email    string    `orm:"unique" hash:"e"`
	Key      string    `orm:"type(text)" hash:"k"`
	CName    string    `orm:"size(32);unique"`
	CEmail   string    `orm:"size(64);unique"`
	Kid      string    `orm:"size(40);unique"`
	Groups   []*Group  `orm:"rel(m2m);rel_through(pygrader-webserver/models.UserGroup)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

func (uu *User) Validate() error {
	if uu == nil {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input")
	} else if uu.Username = strings.Trim(uu.Username, " \t\n"); uu.Username == "" {
		return uti.Errorf(ERR_INVALID_INPUT, "Invalid input 'Username'")
	} else if cname, err := uti.CanonizeName(uu.Username); err != nil {
		return err
	} else {
		uu.CName = cname
	}

	if email, err := uti.CanonizeEmail(uu.Email); err != nil {
		return err
	} else {
		uu.CEmail = email
	}

	if kid, key, err := uti.CryptoGetKeyFingerprint(uu.Key, uti.KeyAll, 40); err != nil {
		return err
	} else {
		uu.Kid = kid
		uu.Key = base64.StdEncoding.EncodeToString(key)
	}
	return nil
}

func AddUser(uu *User) (*User, error) {
	// Compute key fingerprint
	if err := uu.Validate(); err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	if _, err := o.Insert(uu); err != nil {
		return nil, err
	} else {
		return uu, err
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

func GetUsers(email *string, username *string, kid *string) (*[]*User, error) {
	var users []*User
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	cond := orm.NewCondition()
	if email != nil {
		cond = cond.And(`email`, *email)
	}
	if username != nil {
		cond = cond.And(`CName`, *username)
	}
	if kid != nil {
		cond = cond.And(`kid`, *kid)
	}

	if _, err := qs.SetCond(cond).All(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

func UpdateUser(oid int64, uu *User) (*User, error) {
	u := User{Id: oid}
	o := orm.NewOrm()
	// race condition?
	if err := o.Read(&u); err == nil {
		if uu.Key != "" {
			u.Key = uu.Key
		}
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Email != "" {
			u.Email = uu.Email
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
