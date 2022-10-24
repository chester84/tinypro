package models

import (
	"github.com/beego/beego/v2/client/orm"
	"tinypro/common/types"
	"time"
)

func init() {
	orm.RegisterModel(new(Msg))
}

const MSG_TABLENAME = "msg_tab"

type Msg struct {
	ID         int64         `orm:"pk;column(id)" json:"id"` // 主键ID
	UserID     int64         `orm:"column(user_id)" json:"user_id" form:"user_id"`
	CourseId   int64         `json:"course_id" form:"course_id"`
	MsgType    types.MsgEnum `json:"msg_type" form:"msg_type"`
	MsgContent string        `json:"msg_content" form:"msg_content"`
	CreatedAt  time.Time     // 创建时间
}

func (r *Msg) TableName() string {
	return MSG_TABLENAME
}
