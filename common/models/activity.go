package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(Activity))
}

const ACTIVITY_TABLENAME = "activity"

type Activity struct {
	Id        int64  `orm:"pk;" json:"id"`
	Author    string `json:"author"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Bgpic     string `json:"bgpic"`
	CreatedBy int64  `json:"created_by"` // 记录创建者
	CreatedAt int64  `json:"created_at"` // 记录创建时间
	LastOpBy  int64  `json:"last_op_by"` // 最后操作员
	LastOpAt  int64  `json:"last_op_at"` // 最后操作时间
}

func (r *Activity) TableName() string {
	return ACTIVITY_TABLENAME
}
