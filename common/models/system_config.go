package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(SystemConfig))
}

const SYSTEM_CONFIG_TABLENAME string = "system_config"

type SystemConfig struct {
	Id          int64                      `orm:"pk;" json:"id"`
	ItemName    string                     `json:"item_name"`
	Description string                     `json:"description"`
	ItemType    types.SystemConfigItemType `json:"item_type"`
	ItemValue   string                     `json:"item_value"`
	Weight      int                        `json:"weight"`
	Version     int                        `json:"version"`
	Status      types.StatusCommonEnum     `json:"status"`
	OnlineTime  int64                      `json:"online_time"`
	OfflineTime int64                      `json:"offline_time"`
	OpUid       int64                      `orm:"column(op_uid)" json:"op_uid"`
	Ctime       int64                      `json:"ctime"`
	Utime       int64                      `json:"utime"`
}

func (r *SystemConfig) TableName() string {
	return SYSTEM_CONFIG_TABLENAME
}
