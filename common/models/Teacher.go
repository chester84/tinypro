package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(Teacher))
}

const TEACHER_TABLENAME = "teacher_tab"

// Teacher 结构体
type Teacher struct {
	ID        int64     `orm:"pk;column(id)" json:"id"` // 主键ID
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	Name      string    `json:"name" form:"name"`
	Header    string    `json:"header" form:"header"`
}

func (r *Teacher) TableName() string {
	return TEACHER_TABLENAME
}
