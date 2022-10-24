package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(CasbinRbacRole))
}

const CASBIN_RBAC_ROLE_TABLENAME = "casbin_rbac_role"

type CasbinRbacRole struct {
	Id       int64 `orm:"pk;"`
	Name     string
	Status   types.StatusCommonEnum
	LastOpAt int64
	LastOpBy int64
}

func (r *CasbinRbacRole) TableName() string {
	return CASBIN_RBAC_ROLE_TABLENAME
}

func (r *CasbinRbacRole) LoadById(id int64) error {
	o := orm.NewOrm()

	err := o.QueryTable(r.TableName()).Filter("id", id).One(r)
	return err
}

func RoleListValid() (list []CasbinRbacRole) {
	one := CasbinRbacRole{}
	o := orm.NewOrm()

	sql := fmt.Sprintf("SELECT * FROM %s WHERE status = %d ", one.TableName(), types.StatusValid)
	_, err := o.Raw(sql).QueryRows(&list)
	if err != nil {
		logs.Error("[RoleListValid] err:%v", err)
	}
	return
}

func RoleListAll() (list []CasbinRbacRole) {
	one := CasbinRbacRole{}
	o := orm.NewOrm()

	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC", one.TableName())
	_, err := o.Raw(sql).QueryRows(&list)
	if err != nil {
		logs.Error("[RoleListAll] err:%v", err)
	}
	return
}
