package models

import (
	"github.com/beego/beego/v2/client/orm"
	"tinypro/common/types"
	"time"
)

func init() {
	orm.RegisterModel(new(WxInvokeLog))
}

const WX_INVOKE_LOG_TABLENAME = "wx_invoke_log_tab"

type WxInvokeLog struct {
	ID        int64                  `orm:"pk;column(id)" json:"id"` // 主键ID
	ApiType   types.WxInvokeTypeEnum `orm:"column(api_type)" json:"api_type"`
	UserID    int64                  `orm:"column(user_id)" json:"user_id"`
	CourseId  int64                  `orm:"column(course_id)" json:"course_id"`
	Api       string                 `orm:"column(api)" json:"api" `
	Param     string                 `orm:"column(param)" json:"param" `
	RespCode  int                    `orm:"column(resp_code)" json:"resp_code"`
	Resp      string                 `orm:"column(resp)" json:"resp"`
	CreatedAt time.Time              // 创建时间
}

func (r *WxInvokeLog) TableName() string {
	return WX_INVOKE_LOG_TABLENAME
}
