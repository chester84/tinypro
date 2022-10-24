package models

import (
	"encoding/json"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
)

func init() {
	orm.RegisterModel(new(OpLogger))
}

const OP_LOGGER_TABLENAME string = "op_logger"

type OpCodeEnum int

const (
	// op_logger 操作码
	OpCodeAdminUpdate   OpCodeEnum = 100
	OpCodeUpTagsLib     OpCodeEnum = 101
	OpCodeUpAppUser     OpCodeEnum = 102
	OpCodeUpCases       OpCodeEnum = 103
	OpCodeRoleEdit      OpCodeEnum = 104
	OpCodeRuleEdit      OpCodeEnum = 105
	OpCodeUpI18nMapping OpCodeEnum = 106

	OpCodeUpEduOpc    OpCodeEnum = 200
	OpCodeUpEduClass  OpCodeEnum = 201
	OpCodeUpEduCharge OpCodeEnum = 202

	OpCodeUpWatchDevice OpCodeEnum = 300
)

// OpCodeList 描述 opCode 与操作对应关系表
var OpCodeList = map[OpCodeEnum]string{
	OpCodeAdminUpdate: `修改后台用户信息`,
	OpCodeUpTagsLib:   `修改tags-lib`,
	OpCodeUpAppUser:   `修改客户端用户信息`,
	OpCodeUpCases:     `修改Cases`,
	OpCodeRoleEdit:    `rbac角色编辑`,
	OpCodeRuleEdit:    `rbac规则编辑`,
	OpCodeUpEduOpc:    `教育内容运营`,
	OpCodeUpEduClass:  `班级管理`,
	OpCodeUpEduCharge: `费用管理`,

	OpCodeUpWatchDevice: `手表设备`,
}

// OpLogger 描述对应表单行数据结构，及字段映射关系
type OpLogger struct {
	Id        int64      `orm:"pk;" json:"Id,string"`
	OpUid     int64      `orm:"column(op_uid)" json:"OpUid,string"`
	RelatedId int64      `json:"RelatedId,string"`
	OpCode    OpCodeEnum `orm:"column(op_code)"`
	OpTable   string     `orm:"column(op_table)"`
	Original  string
	Edited    string
	Ctime     int64 `json:"Ctime,string"`
}

func (r *OpLogger) TableName() string {
	return OP_LOGGER_TABLENAME
}

// 用于记录一些关键数据被修改的日志,由业务来决定那些需要记录
func OpLogWrite(opUid int64, relatedId int64, opCode OpCodeEnum, opTable string, original interface{}, edited interface{}) {
	originalJson, _ := json.Marshal(original)
	editedJson, _ := json.Marshal(edited)

	opLogIns := OpLogger{
		OpUid:     opUid,
		RelatedId: relatedId,
		OpCode:    opCode,
		OpTable:   opTable,
		Original:  string(originalJson),
		Edited:    string(editedJson),
		Ctime:     libtools.GetUnixMillis(),
	}

	o := orm.NewOrm()
	_, err := o.Insert(&opLogIns)
	if err != nil {
		logs.Error("[OpLogWrite] insert get exception, data: %#v, err: %v", opLogIns, err)
	}
}
