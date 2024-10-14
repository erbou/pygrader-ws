package models

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
	"pygrader-webserver/uti"
)

func init() {
	orm.RegisterModel(new(Group), new(SubGroup))
}

type Group struct {
	Id        int64
	Name      string    `orm:"unique" hash:"n"`
	Scope     *string   `orm:"null" hash:"s"`
	Token     *string   `orm:"null"`
	CName     string    `orm:"unique;size(32);index"`
	Users     []*User   `orm:"reverse(many)"`
	Subgroups []*Group  `orm:"reverse(many)"`
	Supgroups []*Group  `orm:"rel(m2m);rel_through(pygrader-webserver/models.SubGroup)"`
	Created   time.Time `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}

type SubGroup struct {
	Id       int64     `orm:"pk;auto"`
	Subgroup *Group    `orm:"rel(fk);null;on_delete(set_null)"`
	Supgroup *Group    `orm:"rel(fk);null;on_delete(set_null)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`
}

type GroupPreview struct {
	Id    int64
	Name  string
	Scope *string
}

type GroupView struct {
	Id      int64
	Name    string
	Scope   *string
	Token   *string
	Created time.Time
	Updated time.Time
}

// TODO: when supgroup or subgroup are deleted the fields are set to null and stay in the table.
//
//	If both are deleted we accummualte null, null rows because unique does not apply (null != null)
func (obj *SubGroup) TableUnique() [][]string {
	return [][]string{
		{"Supgroup", "Subgroup"},
	}
}

func (obj *SubGroup) TableName() string {
	return "m2m_subgroup"
}

func (obj *Group) Preview() *GroupPreview {
	if obj == nil {
		return nil
	}

	return &GroupPreview{Id: obj.Id, Name: obj.Name, Scope: obj.Scope}
}

func (obj *Group) View() *GroupView {
	if obj == nil {
		return nil
	}

	return &GroupView{
		Id:      obj.Id,
		Name:    obj.Name,
		Scope:   obj.Scope,
		Token:   obj.Token,
		Created: obj.Created,
		Updated: obj.Updated,
	}
}

func addSubGroup(g *Group, subg *Group) (*SubGroup, error) {
	o := orm.NewOrm()
	gG := SubGroup{
		Supgroup: g,
		Subgroup: subg,
		Created:  time.Now().UTC(),
		Updated:  time.Now().UTC(),
	}

	if nrow, err := o.InsertOrUpdate(&gG); nrow > 0 {
		return &gG, err
	} else {
		err := o.Raw("SELECT T0.`id` FROM `m2m_subgroup` T0 WHERE T0.`supgroup_id` = ? AND T0.`subgroup_id` = ?",
			g.Id, subg.Id).QueryRow(&gG)

		return &gG, err
	}
}

func loadRelatedSup(gid int64) *[]*GroupPreview {
	var relations []*GroupPreview

	o := orm.NewOrm()
	if _, err := o.Raw("SELECT T0.`id`, T0.`name`, T0.`scope` FROM `group` T0 INNER JOIN `m2m_subgroup` T1 ON T1.`supgroup_id` = T0.`id` AND  T1.`subgroup_id` = ?", gid).QueryRows(&relations); err != nil {
		logs.Warning("%v supgroup=%v in m2m_subgroup", err, gid)
	}

	return &relations
}

func loadRelatedSub(gid int64) *[]*GroupPreview {
	var relations []*GroupPreview

	// load related sub without tokens
	o := orm.NewOrm()
	if _, err := o.Raw("SELECT T0.`id`, T0.`name`, T0.`scope` FROM `group` T0 INNER  JOIN `m2m_subgroup` T1 ON T1.`subgroup_id` = T0.`id` AND  T1.`supgroup_id` = ?", gid).QueryRows(&relations); err != nil {
		logs.Warning("%v subgroup=%v in m2m_subgroup", err, gid)
	}

	return &relations
}

func (obj *Group) Validate() error {
	if obj == nil {
		return uti.Errorf(ErrInvalidInput, "Invalid Input")
	} else if obj.Name = strings.Trim(obj.Name, " \t\n"); obj.Name == "" {
		return uti.Errorf(ErrInvalidInput, "Invalid Input")
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

	if _, err := o.Insert(obj); err == nil {
		return obj, nil
	} else {
		return nil, err
	}
}

func GetGroup(oid int64) (*Group, error) {
	g := Group{Id: oid}
	o := orm.NewOrm()

	if err := o.Read(&g); err != nil {
		return nil, err
	} else {
		return &g, nil
	}
}

func GetGroups(name *string, page int, pageSize int) (*[]*GroupPreview, error) {
	var groups []*Group

	o := orm.NewOrm()
	qs := o.QueryTable("group")
	cond := orm.NewCondition()

	if name != nil {
		if _name, err := uti.CanonizeName(*name); err != nil {
			return nil, err
		} else {
			cond = cond.And(`name`, _name)
		}
	}

	if _, err := qs.Limit(pageSize, (page-1)*pageSize).SetCond(cond).All(&groups, `Id`, `Name`, `Scope`); err != nil {
		return nil, err
	}

	list := make([]*GroupPreview, 0, len(groups))
	for _, g := range groups {
		list = append(list, g.Preview())
	}

	return &list, nil
}

func UpdateGroup(oid int64, obj *Group) (*Group, error) {
	if obj == nil {
		return nil, uti.Errorf(ErrInvalidInput, "Invalid Input")
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

		if _, err := o.Update(&dbObj); err != nil {
			return nil, err
		}

		return &dbObj, nil
	}
}

func DeleteGroup(oid int64) (int64, error) {
	o := orm.NewOrm()

	return o.Delete(&Group{Id: oid})
}

func GroupAddUser(g *Group, u *User, acl string) (int64, error) {
	uG, err := AddUserGroup(u, g, acl)

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

func GroupRemoveUser(gid int64, uid int64) (int64, error) {
	o := orm.NewOrm()
	// QueryM2M.Remove does not remove the row in m2m relationships
	if res, err := o.Raw("DELETE FROM `m2m_user_group` WHERE `group_id` = ? AND `user_id` = ?", gid, uid).Exec(); err != nil {
		return 0, err
	} else if n, err := res.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func GroupAddSubgroup(g *Group, subg *Group) (int64, error) {
	if g.Id == subg.Id {
		return 0, uti.Errorf(ErrInvalidInput, "Invalid Input")
	}

	gG, err := addSubGroup(g, subg)

	return gG.Id, err
}

func GroupGetSubgroups(oid int64) (*[]*GroupPreview, error) {
	list := loadRelatedSub(oid)

	return list, nil
}

func GroupGetSupgroups(oid int64) (*[]*GroupPreview, error) {
	list := loadRelatedSup(oid)

	return list, nil
}

func GroupRemoveSubgroup(gid int64, sgid int64) (int64, error) {
	o := orm.NewOrm()
	// QueryM2M.Remove does not remove the row in m2m relationships
	if res, err := o.Raw("DELETE FROM `m2m_subgroup` WHERE `supgroup_id` = ? AND `subgroup_id` = ?", gid, sgid).Exec(); err != nil {
		return 0, err
	} else if n, err := res.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}
