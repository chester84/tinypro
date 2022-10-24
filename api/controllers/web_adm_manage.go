package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/cerror"
	"tinypro/common/lib/device"
	"tinypro/common/models"
	"tinypro/common/pkg/admin"
	"tinypro/common/pkg/helper"
	"tinypro/common/pkg/rbacv2"
	"tinypro/common/pkg/system/config"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebAdmManageController struct {
	WebAdmBaseController
}

func (c *WebAdmManageController) Prepare() {
	// 调用上一级的 Prepare 方
	c.WebAdmBaseController.Prepare()
}

func (c *WebAdmManageController) AdmUserList() {
	var admList []models.Admin
	admObj := models.Admin{}
	condBox := map[string]interface{}{}

	total, _ := models.OrmList(&admObj, condBox, 1, 500, true, &admList)

	type admin4Api struct {
		Id            string                `json:"id"`
		Email         string                `json:"email"`
		Nickname      string                `json:"nickname"`
		Roles         []string              `json:"roles"`
		Status        types.StatusBlockEnum `json:"status"`
		RegisterTime  string                `json:"registerTime"`
		LastLoginTime string                `json:"lastLoginTime"`
	}
	var list []admin4Api
	for _, adm := range admList {
		roleNameBox, _ := rbacv2.OperatorAccessRole(adm.Id)

		list = append(list, admin4Api{
			Id:            libtools.Int642Str(adm.Id),
			Email:         adm.Email,
			Nickname:      adm.Nickname,
			Roles:         roleNameBox,
			Status:        adm.Status,
			RegisterTime:  libtools.Int642Str(adm.RegisterTime / 1000),
			LastLoginTime: libtools.Int642Str(adm.LastLoginTime / 1000),
		})
	}

	c.SuccessResponse(types.H{
		"list":  list,
		"total": total,
	})
}

func (c *WebAdmManageController) AdmUserSave() {
	var affectedId = ""

	email := c.GetString("email")
	nickname := c.GetString("nickname")
	roles := c.GetStrings("roles")
	password := c.GetString("password")

	opAction := c.GetString("opAction")

	if len(email) < 2 || len(nickname) < 2 {
		logs.Warning("参数非法")
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	if opAction == "create" {

		if len(password) < 6 {
			logs.Warning("密码为空")
			c.TerminateWithCode(cerror.ParameterValueOutOfRange)
			return
		}

		registerTime := libtools.GetUnixMillis()
		adminUID, _ := device.GenerateBizId(types.AccountSystem)
		adminModel := models.Admin{
			Id:           adminUID,
			Email:        email,
			Nickname:     nickname,
			Password:     libtools.PasswordEncrypt(password, registerTime),
			RegisterTime: registerTime,
			Status:       types.Unblock,
		}
		_, err := admin.Add(&adminModel)
		if err != nil {
			logs.Error("[AdmUserSave] db insert exception, data: %#v, err: %v", adminModel, err)
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}

		for _, role := range roles {
			roleObj, err := rbacv2.RoleByPkName(role)
			if err != nil {
				logs.Info("role: %s, err: %v", role, err)
				continue
			}

			_, _ = rbacv2.Enforcer().AddNamedGroupingPolicy("g", libtools.Int642Str(adminUID), libtools.Int642Str(roleObj.Id))
		}

		affectedId = libtools.Int642Str(adminUID)
	} else {

		var cols []string
		id, _ := c.GetInt64("id")
		modelData, err := admin.OneByUid(id)
		if err != nil {
			c.TerminateWithCode(cerror.InvalidData)
			return
		}

		originData := modelData

		modelData.Nickname = nickname
		cols = append(cols, "nickname")

		modelData.Email = email
		cols = append(cols, "email")

		if password != "" {
			//不空的情况下再更新
			modelData.Password = libtools.PasswordEncrypt(password, originData.RegisterTime)
			if originData.Password != modelData.Password {
				cols = append(cols, "Password")
			}
		}

		_, err = admin.Update(&modelData, &originData, cols)
		if err != nil {
			logs.Error("db update exception, data: %#v, err: %v", modelData, err)
			return
		}

		// save policy
		orgRole, _ := rbacv2.Enforcer().GetRolesForUser(libtools.Int642Str(modelData.Id))
		for _, r := range orgRole {
			_, _ = rbacv2.Enforcer().RemoveNamedGroupingPolicy("g", libtools.Int642Str(modelData.Id), r)
		}

		for _, role := range roles {
			roleObj, err := rbacv2.RoleByPkName(role)
			if err != nil {
				logs.Info("role: %s, err: %v", role, err)
				continue
			}

			_, _ = rbacv2.Enforcer().AddNamedGroupingPolicy("g", libtools.Int642Str(modelData.Id), libtools.Int642Str(roleObj.Id))
		}

		// 写操作日志
		models.OpLogWrite(c.AccountId, modelData.Id, models.OpCodeAdminUpdate, modelData.TableName(), originData, modelData)
		affectedId = libtools.Int642Str(modelData.Id)
	}

	_ = rbacv2.Enforcer().SavePolicy()
	rbacv2.RefreshRbacCache()

	c.SuccessResponse(types.H{
		"affectedId": affectedId,
	})
}

func (c *WebAdmManageController) AdmUserChangeStatus() {
	id, _ := c.GetInt64("id")
	if id <= types.SuperAdminUID {
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	opAction := c.GetString(`opAction`)

	var adm models.Admin
	err := models.OrmOneByPkId(id, &adm)
	if err != nil {
		c.TerminateWithCode(cerror.InvalidData)
		return
	}

	switch opAction {
	case "block":
		adm.Status = types.Block

	case "unblock":
		adm.Status = types.Unblock
	}

	_, err = models.OrmUpdate(&adm, []string{"Status"})
	if err != nil {
		logs.Warning("db update exception, adminId: %d, data: %#v, err: %v", c.AccountId, adm, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(types.H{
		"status": adm.Status,
	})
}

func (c *WebAdmManageController) RbacRuleList() {
	list := rbacv2.RuleG2List()
	c.SuccessResponse(types.H{
		"list": list,
	})
}

func (c *WebAdmManageController) RbacRuleLSave() {
	var affectedId string

	opAction := c.GetString("op_action")
	if opAction != "create" && opAction != "edit" {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	rule := c.GetString("v1")
	router := c.GetString("v0")
	if rule == "" || router == "" {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	if opAction == "edit" {
		id, _ := c.GetInt64("Id")
		if id <= 0 {
			c.TerminateWithCode(cerror.LostRequiredParameters)
			return
		}
		one, err := rbacv2.RuleByPkID(id)
		if err != nil {
			c.TerminateWithCode(cerror.LostRequiredParameters)
			return
		}

		origin := one

		one.V0 = router
		one.V1 = rule

		_, err = models.OrmUpdate(&one, []string{"v0", "v1"})
		if err != nil {
			logs.Warning(fmt.Sprintf("修改失败, OrmUpdate one: %#v, err: %v", one, err))
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}

		models.OpLogWrite(c.AccountId, one.Id, models.OpCodeRuleEdit, one.TableName(), origin, one)
		affectedId = libtools.Int642Str(one.Id)
	} else {
		if rbacv2.RuleCheckByRouter(router) {
			logs.Warning(fmt.Sprintf("规则/路由 已存在, rule: %s", router))
			c.TerminateWithCode(cerror.RepeatedSubmitData)
			return
		}

		one := models.CasbinRule{
			PType: "g2",
			V0:    router,
			V1:    rule,
		}
		lastID, err := models.OrmInsert(&one)
		if err != nil {
			logs.Warning(fmt.Sprintf("新增失败, OrmInsert id:%d err: %v", lastID, err))
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}

		affectedId = libtools.Int642Str(lastID)
	}

	c.SuccessResponse(types.H{
		"affectedId": affectedId,
	})
}

func (c *WebAdmManageController) RbacPolicyConfig() {
	c.SuccessResponse(types.H{
		"config": rbacv2.Enforcer().GetAllNamedRoles("g2"),
	})
}

func (c *WebAdmManageController) RbacRoleConfig() {
	var config = make([]string, 0)
	allRole := models.RoleListAll()
	for _, one := range allRole {
		config = append(config, one.Name)
	}

	c.SuccessResponse(types.H{
		"config": config,
	})
}

func (c *WebAdmManageController) RbacRoleList() {
	type listItem struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		OwnPolicies []string `json:"ownPolicies"`
	}

	allRole := models.RoleListAll()
	var list []listItem
	for _, one := range allRole {
		list = append(list, listItem{
			Id:          libtools.Int642Str(one.Id),
			Name:        one.Name,
			OwnPolicies: rbacv2.RolePolicy(libtools.Int642Str(one.Id)),
		})
	}

	c.SuccessResponse(types.H{
		"list": list,
	})
}

func (c *WebAdmManageController) RbacRoleSave() {
	name := c.GetString("name")
	ownPolicies := c.GetStrings("ownPolicies")
	id, _ := libtools.Str2Int64(c.GetString("id"))
	now := libtools.GetUnixMillis()

	one := models.CasbinRbacRole{}
	origin := models.CasbinRbacRole{}
	if id > 0 {
		err := one.LoadById(id)
		if err != nil {
			logs.Warning(fmt.Sprintf("修改失败,LoadById id:%d err: %v", id, err))
			c.TerminateWithCode(cerror.IncompleteData)
			return
		}
		origin = one
		one.Name = name
		one.LastOpBy = c.AccountId
		one.LastOpAt = now
		_, err = models.OrmUpdate(&one, []string{"Name", "LastOpBy", "LastOpAt"})
		if err != nil {
			logs.Warning(fmt.Sprintf("修改失败,OrmUpdate id:%d err: %v", id, err))
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}
	} else {
		one.Status = types.StatusValid
		one.Name = name
		one.LastOpBy = c.AccountId
		one.LastOpAt = now
		id, err := models.OrmInsert(&one)
		if err != nil {
			logs.Warning(fmt.Sprintf("新增失败, OrmInsert id:%d err: %v", id, err))
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}
		one.Id = id
	}

	models.OpLogWrite(c.AccountId, one.Id, models.OpCodeRoleEdit, one.TableName(), origin, one)

	// save policy
	_, err := rbacv2.Enforcer().RemoveFilteredNamedPolicy("p", 0, libtools.Int642Str(one.Id))
	if err != nil {
		logs.Error("[RoleSave] rbacv2.Enforcer().RemoveFilteredNamedPolicy exception, err: %v", err)
	}

	for _, policy := range ownPolicies {
		_, err = rbacv2.Enforcer().AddPolicy(libtools.Int642Str(one.Id), policy, rbacv2.PolicyStr)
		if err != nil {
			logs.Error("[RoleSave] AddPolicy get exception, err: %v", err)
		}
	}
	err = rbacv2.Enforcer().SavePolicy()
	if err != nil {
		logs.Error("[RoleSave] get exception, err: %v", err)
	}

	rbacv2.RefreshRbacCache()

	c.SuccessResponse(types.H{
		"affectedId": libtools.Int642Str(one.Id),
	})
}

func (c *WebAdmManageController) OpLog() {
	condBox := map[string]interface{}{}
	c.ParseDateRangeCommon(`ctime`, condBox, true)

	opCode, _ := c.GetInt("opCode")
	if opCode > 0 {
		condBox["opCode"] = opCode
	}

	id, _ := c.GetInt64("id")
	if id > 0 {
		condBox["id"] = id
	}

	relatedId, _ := c.GetInt64("relatedId")
	if relatedId > 0 {
		condBox["relatedId"] = relatedId
	}

	opUID, err := c.GetInt64("opUserId")
	if err == nil && opUID > -1 {
		condBox["opUid"] = opUID
		c.Data["opUid"] = opUID
	} else {
		c.Data["opUid"] = -1
	}

	// 分页逻辑
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")

	var (
		logObj  models.OpLogger
		listBox []models.OpLogger
	)

	type itemT struct {
		models.OpLogger
		OpUserName string
		OpCodeDesc string
	}

	total, _ := models.OrmList(&logObj, condBox, page, limit, true, &listBox)

	var list = make([]itemT, 0)
	for _, obj := range listBox {
		list = append(list, itemT{
			OpLogger:   obj,
			OpUserName: helper.OperatorName(obj.OpUid),
			OpCodeDesc: models.OpCodeList[obj.OpCode],
		})
	}

	c.SuccessResponse(types.H{
		"list":  list,
		"total": total,
	})
}

func (c *WebAdmManageController) SystemConfigList() {
	condBox := map[string]interface{}{}

	status := c.GetString("status")
	if status != "" {
		isValid, _ := libtools.Str2Int(status)
		condBox["status"] = isValid
	}
	itemName := c.GetString("item_name")
	if len(itemName) > 0 {
		condBox["item_name"] = itemName
	}

	conf := types.SystemConfigItemTypeMap()
	listBox, _, _ := config.List(condBox)

	type itemT struct {
		models.SystemConfig

		OpUserName   string `json:"op_user_name"`
		ItemTypeDesc string `json:"item_type_desc"`
	}
	var list = make([]itemT, 0)
	for _, obj := range listBox {
		list = append(list, itemT{
			SystemConfig: obj,

			OpUserName:   helper.OperatorName(obj.OpUid),
			ItemTypeDesc: conf[obj.ItemType],
		})
	}

	c.SuccessResponse(types.H{"list": list, "typeMap": conf})
}

func (c *WebAdmManageController) SystemConfigSave() {
	itemName := c.GetString("item_name")
	itemTypeP := c.GetString("item_type")
	itemTypeInt, _ := libtools.Str2Int(itemTypeP)
	itemType := types.SystemConfigItemType(itemTypeInt)
	itemValue := c.GetString("item_value")
	weight, _ := c.GetInt("weight")
	description := c.GetString("description")

	logs.Debug("111 itemType %d", itemType)

	if len(itemName) <= 0 || len(itemValue) <= 0 || itemType <= 0 {
		logs.Debug("itemName %s", itemName)
		logs.Debug("itemValue %s", itemValue)
		logs.Debug("itemType %d", itemType)
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	_, err := config.Create(itemName, itemValue, itemType, weight, description, c.AccountId)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(``)
}

func (c *WebAdmManageController) TagsList() {
	condBox := map[string]interface{}{}

	status := c.GetString("status")
	if status != "" {
		isValid, _ := libtools.Str2Int(status)
		condBox["status"] = isValid
	}
	name := c.GetString("name")
	if len(name) > 0 {
		condBox["name"] = name
	}

	var (
		obj     models.TagsLibrary
		objList []models.TagsLibrary
	)

	_, _ = models.OrmList(&obj, condBox, 1, 500, false, &objList)

	type itemT struct {
		models.TagsLibrary

		SN         string `json:"sn"`
		OpUserName string
	}
	var list = make([]itemT, 0)
	for _, obj := range objList {
		list = append(list, itemT{
			TagsLibrary: obj,

			OpUserName: helper.OperatorName(obj.LastOpBy),
			SN:         libtools.Int642Str(obj.Id),
		})
	}

	c.SuccessResponse(types.H{"list": list})
}

func (c *WebAdmManageController) TagsSave() {
	var err error
	var one models.TagsLibrary

	opAction := c.GetString(`opAction`)

	name := c.GetString(`Name`)
	detail := c.GetString(`Detail`)
	weight, weightErr := c.GetInt(`Weight`)
	timeNow := libtools.GetUnixMillis()

	if opAction == "create" {
		pkId, _ := device.GenerateBizId(types.TagsLibBiz)

		one = models.TagsLibrary{
			Id:        pkId,
			Name:      name,
			Detail:    detail,
			Weight:    weight,
			Status:    types.StatusValid,
			CreatedAt: timeNow,
			CreatedBy: c.AccountId,
			LastOpAt:  timeNow,
			LastOpBy:  c.AccountId,
		}

		_, err = models.OrmInsert(&one)
		if err != nil {
			logs.Error("[TagsSave] db insert exception, ip: %s, opUid: %d, one: %#v, err: %v", c.IP, c.AccountId, one, err)
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}
	} else {
		id, _ := c.GetInt64(`id`)
		err = models.OrmOneByPkId(id, &one)
		if err != nil {
			logs.Error("[TagsSave] db get origin data exception, ip: %s, opUid: %d, id: %d, err: %v",
				c.IP, c.AccountId, id, err)
			c.TerminateWithCode(cerror.InvalidData)
			return
		}

		origin := one

		if name != "" {
			one.Name = name
		}
		if detail != "" {
			one.Detail = detail
		}
		if weightErr == nil {
			one.Weight = weight
		}

		one.LastOpBy = c.AccountId
		one.LastOpAt = libtools.GetUnixMillis()

		switch c.GetString(`opAction`) {
		case "online":
			one.Status = types.StatusValid

		case "offline":
			one.Status = types.StatusInvalid
		}

		_, err = models.OrmAllUpdate(&one)
		if err != nil {
			logs.Error("[TagsSave] db update exception, ip: %s, opUid: %d, one: %#v, err: %v", c.IP, c.AccountId, one, err)
			c.TerminateWithCode(cerror.ServiceDbOpFail)
			return
		}

		models.OpLogWrite(c.AccountId, one.Id, models.OpCodeUpTagsLib, one.TableName(), origin, one)
	}

	c.SuccessResponse(types.H{
		"affected_sn": libtools.Int642Str(one.Id),
	})
}
