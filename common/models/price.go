package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(Price))
}

const PRICE_TABLENAME = "price"

type Price struct {
	Id          int64  `orm:"pk;" json:"id"`
	CourseId    int64  `json:"course_id,string"`
	Price       int64  `json:"price"`        //
	OriginPrice int64  `json:"origin_price"` //
	DateBegin   string `json:"date_begin"`
	DateEnd     string `json:"date_end"`
	CreatedBy   int64  `json:"created_by"` // 记录创建者
	CreatedAt   int64  `json:"created_at"` // 记录创建时间
	LastOpBy    int64  `json:"last_op_by"` // 最后操作员
	LastOpAt    int64  `json:"last_op_at"` // 最后操作时间
}

func (r *Price) TableName() string {
	return PRICE_TABLENAME
}
