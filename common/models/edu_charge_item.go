package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(EduChargeItem))
}

const EDU_CHARGE_ITEM_TABLENAME = "edu_charge_item"

type EduChargeItem struct {
	Id           int64 `orm:"pk;" json:"id,string"`
	ChargeConfId int64 `json:"charge_conf_id,string"` // 关联的费用配置id
	Name         string
	ChargeType   types.EduChargeTypeEnum
	Amount       int64
	Status       types.StatusCommonEnum
	CreatedBy    int64 `json:"created_by,string"` // 记录创建者
	CreatedAt    int64 // 记录创建时间
	LastOpBy     int64 `json:"last_op_by,string"` // 最后操作员
	LastOpAt     int64 // 最后操作时间
}

func (r *EduChargeItem) TableName() string {
	return EDU_CHARGE_ITEM_TABLENAME
}

type EduChargeItemExpand struct {
	EduChargeItem
	AmountHuman string
	CreatedUser string
	LastOpUser  string
}
