package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(EduChargeConf))
}

const EDU_CHARGE_CONF_TABLENAME = "edu_charge_conf"

type EduChargeConf struct {
	Id                int64 `orm:"pk;" json:"id,string"`
	ClassId           int64
	FixedAmountSnap   string
	UnfixedAmountSnap string
	FixedAmount       int64 // 固定总金额
	Status            types.StatusCommonEnum
	CreatedBy         int64 // 记录创建者
	CreatedAt         int64 // 记录创建时间
	LastOpBy          int64 // 最后操作员
	LastOpAt          int64 // 最后操作时间
}

func (r *EduChargeConf) TableName() string {
	return EDU_CHARGE_CONF_TABLENAME
}

type EduChargeConfExpand struct {
	EduChargeConf
	ClassName   string
	CreatedUser string
	LastOpUser  string
}
