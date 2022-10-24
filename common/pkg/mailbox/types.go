package mailbox

type ListItem struct {
	Id         int64 `orm:"pk;"`
	MailID     int64 `orm:"column(mail_id)"`
	SenderID   int64 `orm:"column(sender_id)"`
	ReceiverID int64 `orm:"column(receiver_id)"`
	Title      string
	Body       string
	ReceiveAt  int64
	ReadAt     int64
	MakeReadAt int64
	DeleteAt   int64
}
