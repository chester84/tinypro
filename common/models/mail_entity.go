package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(new(MailEntity))
}

const MAIL_ENTITY_TABLENAME = "mail_entity"

type MailEntity struct {
	Id       int64 `orm:"pk;"`
	Owner    int64
	Title    string
	Body     string
	Ctime    int64
	DeleteAt int64
}

func (r *MailEntity) TableName() string {
	return MAIL_ENTITY_TABLENAME
}
