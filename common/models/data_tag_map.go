package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(DataTagMap))
}

const DATA_TAG_MAP_TABLENAME = "data_tag_map"

type DataTagMap struct {
	Id       int64                  `orm:"pk;"`
	DataID   int64                  `orm:"column(data_id)"` // 数据主键ID
	TagID    int64                  `orm:"column(tag_id)"`  // 标签ID
	BizSN    types.BizSN            `orm:"column(biz_sn)"`  // 数据所属业务
	Status   types.StatusCommonEnum // 0: 无效; 1: 有效
	LastOpBy int64                  // 最后操作员
	LastOpAt int64                  // 最后操作时间
}

func (r *DataTagMap) TableName() string {
	return DATA_TAG_MAP_TABLENAME
}
