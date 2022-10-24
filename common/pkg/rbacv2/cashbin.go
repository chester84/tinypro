package rbacv2

import (
	"github.com/beego/beego/v2/core/logs"
	c2 "github.com/casbin/casbin/v2"
)

const confPath = "conf/rbac_model.conf"

var enforcer *c2.Enforcer

func init() {

	//var err error
	//a := NewAdapter()
	//enforcer, err = c2.NewEnforcer(confPath, a)
	//if err != nil {
	//	logs.Error("[rbac2] NewEnforcer err:%v", err)
	//	os.Exit(1)
	//}
	//err = enforcer.LoadPolicy()
	//if err != nil {
	//	logs.Error("[rbac2] LoadPolicy err:%v", err)
	//	os.Exit(1)
	//}
}

func Enforcer() *c2.Enforcer {
	// todo 优化不要每次都读数据库
	a := NewAdapter()
	e, err := c2.NewEnforcer(confPath, a)
	if err != nil {
		logs.Error("[Enforcer] NewEnforcer err:%v", err)
		return enforcer
	}

	err = e.LoadPolicy()
	if err != nil {
		logs.Error("[Enforcer] LoadPolicy err:%v", err)
		return enforcer
	}
	return e
}

func SavePolicy(enforcer *c2.Enforcer) {
	//todo lock the save modify
	err := enforcer.SavePolicy()
	if err != nil {
		logs.Error("[SavePolicy] get exception, err: %v", err)
	}
}

func RolePolicy(roleIds string) []string {
	// ps  [[编辑角色 创建模板权限 write] [编辑角色 发布文章权限 write]]
	ps := Enforcer().GetFilteredPolicy(0, roleIds)
	mps := []string{}
	for _, v := range ps {
		if len(v) > 2 {
			mps = append(mps, v[1])
		}
	}
	return mps
}
