package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(new(AdminLoginLog))
}

const ADMIN_LOGIN_LOG_TABLENAME string = "admin_login_log"

type AdminLoginLog struct {
	Id       int64  `orm:"pk;"`
	AdminUID int64  `orm:"column(admin_uid)"`
	IP       string `orm:"column(ip)"`
	Ctime    int64
}

func (*AdminLoginLog) TableName() string {
	return ADMIN_LOGIN_LOG_TABLENAME
}
