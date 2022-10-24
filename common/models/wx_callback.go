package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(WxCallback))
}

const WX_CALLBACK_TABLENAME = "wx_callback"

type WxCallback struct {
	Id        int64 `orm:"pk;"`
	ReqType   types.WxCallbackReqTypeEnum
	ReqUrl    string
	ReqParams string
	RespCode  int
	Resp      string
	CreatedAt int64
}

func (r *WxCallback) TableName() string {
	return WX_CALLBACK_TABLENAME
}
