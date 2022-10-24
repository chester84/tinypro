package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(new(WxOpenId))
}

const WX_OPEN_ID_TABLENAME = "wx_open_id"

type WxOpenId struct {
	Id        int64 `orm:"pk;"`
	UserId    int64
	AppSN     int `orm:"column(app_sn)"`
	OpenId    string
	CreatedAt int64
	LastOpAt  int64
}

func (r *WxOpenId) TableName() string {
	return WX_OPEN_ID_TABLENAME
}
