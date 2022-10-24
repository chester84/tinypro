package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(new(Mailbox))
}

const MAILBOX_TABLENAME = "mailbox"

type Mailbox struct {
	Id         int64 `orm:"pk;"`
	MailID     int64 `orm:"column(mail_id)"`
	SenderID   int64 `orm:"column(sender_id)"`
	ReceiverID int64 `orm:"column(receiver_id)"`
	ReceiveAt  int64
	ReadAt     int64
	MakeReadAt int64
	DeleteAt   int64
}

func (r *Mailbox) TableName() string {
	return MAILBOX_TABLENAME
}
