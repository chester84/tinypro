package rbacv2

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/cache"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func HasPolicyByIdRouter(idStr string, router string) bool {
	if idStr == "1" {
		return true
	}

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	// 从缓存读取
	filed := idStr + "_" + router + "_" + PolicyStr
	if v, err := redis.Bool(cacheClient.Do("HGET", policyByIdRouter, filed)); err == nil {
		return v
	}

	b, err := Enforcer().Enforce(idStr, router, PolicyStr)
	if err != nil {
		logs.Error("[HasPolicyByIdRouter] Enforce err:%v idStr:%v router:%v", err, idStr, router)
		return false
	}

	//写入缓存
	_, err = cacheClient.Do("HSET", policyByIdRouter, filed, b)
	if err != nil {
		logs.Error("[HasPolicyByIdRouter] redis> HSET %s %s %v", policyByIdRouter, filed, b)
	}

	return b
}

func HasRoleForUserByUid(uId int64, roleId int64) bool {
	if uId == types.SuperAdminUID {
		return true
	}

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	// 从缓存读取
	filed := libtools.Int642Str(uId) + "_" + libtools.Int642Str(roleId)
	if v, err := redis.Bool(cacheClient.Do("HGET", policyByIdRouter, filed)); err == nil {
		return v
	}

	rv, err := Enforcer().HasRoleForUser(libtools.Int642Str(uId), libtools.Int642Str(roleId))
	if err != nil {
		logs.Error("[HasRoleForUserByUid] uId:%d roleId:%d", uId, roleId)
		return false
	}
	//写入缓存
	_, err = cacheClient.Do("HSET", policyByIdRouter, filed, rv)
	if err != nil {
		logs.Error("[HasRoleForUserByUid] redis> HSET %s %s %v", policyByIdRouter, filed, rv)
	}

	return rv
}

func RefreshRbacCache() {
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	_, err := cacheClient.Do("DEL", policyByIdRouter)
	if err != nil {
		logs.Error("[RefreshRbacCache] redis> DEL %s", policyByIdRouter)
	}
}

func RuleG2List() []models.CasbinRule {
	var list = make([]models.CasbinRule, 0)

	o := orm.NewOrm()
	m := models.CasbinRule{}

	_, err := o.QueryTable(m.TableName()).Filter("p_type", "g2").OrderBy("v1").All(&list)
	if err != nil {
		logs.Error("[RuleG2List] db get exception, err: %v", err)
	}

	return list
}

func RoleByPkID(id int64) (models.CasbinRbacRole, error) {
	m := models.CasbinRbacRole{}

	if id <= 0 {
		err := fmt.Errorf(`input id is 0`)
		logs.Warning("[RoleByPkID] err: %v", err)
		return m, err
	}

	o := orm.NewOrm()

	err := o.QueryTable(m.TableName()).Filter("id", id).One(&m)
	if err != nil {
		if err.Error() != types.EmptyOrmStr {
			logs.Error("[RoleByPkID] db exception, err: %v", err)
		}
	}

	return m, err
}

func RoleByPkName(name string) (models.CasbinRbacRole, error) {
	m := models.CasbinRbacRole{}

	if name == "" {
		err := fmt.Errorf(`input id is empty`)
		logs.Warning("[RoleByPkName] err: %v", err)
		return m, err
	}

	o := orm.NewOrm()

	err := o.QueryTable(m.TableName()).Filter("name", name).One(&m)
	if err != nil {
		if err.Error() != types.EmptyOrmStr {
			logs.Error("[RoleByPkName] db exception, err: %v", err)
		}
	}

	return m, err
}

func RuleByPkID(id int64) (models.CasbinRule, error) {
	m := models.CasbinRule{}

	if id <= 0 {
		err := fmt.Errorf(`input id is 0`)
		logs.Warning("[RuleByPkID] err: %v", err)
		return m, err
	}

	o := orm.NewOrm()

	err := o.QueryTable(m.TableName()).Filter("id", id).One(&m)
	if err != nil {
		if err.Error() != types.EmptyOrmStr {
			logs.Error("[RuleByPkID] db exception, err: %v", err)
		}
	}

	return m, err
}

func RuleCheckByRouter(router string) bool {
	m := models.CasbinRule{}
	o := orm.NewOrm()

	err := o.QueryTable(m.TableName()).Filter("v0", router).One(&m)
	if err != nil {
		if err.Error() != types.EmptyOrmStr {
			logs.Error("[RuleCheckByRouter] db exception, err: %v", err)
		}
		return false
	} else {
		return true
	}
}

func HasAccessWithRoleName(uid int64, roleName ...string) bool {
	if uid == types.SuperAdminUID {
		return true
	}

	var checkBox = map[string]bool{}
	userRoleList, _ := Enforcer().GetRolesForUser(libtools.Int642Str(uid))
	for _, rid := range userRoleList {
		roleID, _ := libtools.Str2Int64(rid)
		one, err := RoleByPkID(roleID)
		if err != nil {
			continue
		}

		checkBox[one.Name] = true
	}

	for _, role := range roleName {
		if checkBox[role] {
			return true
		}
	}

	return false
}

func OperatorAccessRole(opId int64) (roleNameBox []string, roleBox []int64) {
	roleNameBox = []string{}
	roleBox = []int64{}

	if opId == types.SuperAdminUID {
		roleNameBox = append(roleNameBox, "super-admin")
	}

	userRoleList, _ := Enforcer().GetRolesForUser(libtools.Int642Str(opId))
	//logs.Debug("userRoleList: %#v", userRoleList)

	for _, rid := range userRoleList {
		roleId, _ := libtools.Str2Int64(rid)
		one, err := RoleByPkID(roleId)
		if err != nil {
			continue
		}

		roleNameBox = append(roleNameBox, one.Name)
		roleBox = append(roleBox, roleId)
	}

	if len(roleNameBox) == 0 {
		// 给一个空权限,防止前端卡住
		roleNameBox = append(roleNameBox, "null")
	}

	return
}
