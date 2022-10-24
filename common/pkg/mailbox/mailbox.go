package mailbox

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/device"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func OneByPKID(pkID int64) (one models.Mailbox, err error) {
	one.Id = pkID

	o := orm.NewOrm()

	err = o.Read(&one)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[mailbox->OneByPKID] read one get exception, id: %d, err: %v", pkID, err)
	}

	return
}

func OneItemByPKID(pkID int64) (item ListItem, err error) {
	one := models.Mailbox{}
	one.Id = pkID

	o := orm.NewOrm()

	err = o.Read(&one)
	if err != nil {
		logs.Warning("[mailbox] OneByPKID mailbox data does not exist, id: %d", pkID)
		return
	}

	mailEntity := models.MailEntity{
		Id: one.MailID,
	}
	err = o.Read(&mailEntity)
	if err != nil {
		logs.Warning("[mailbox] OneByPKID mail entity data does not exist, id: %d", pkID)
		return
	}

	item.Id = one.Id
	item.ReceiverID = one.ReceiverID
	item.Title = mailEntity.Title
	item.Body = mailEntity.Body
	item.ReceiveAt = one.ReceiveAt

	return
}

func SendMail(senderID int64, receiverIDBox []int64, title, body string) (mailID int64, mailboxList []int64, err error) {
	if len(receiverIDBox) == 0 {
		err = fmt.Errorf(`no receiver`)
		logs.Warning("[SendMail] err: %v", err)
		return
	}
	if title == "" || body == "" {
		err = fmt.Errorf(`empty title or body`)
		logs.Warning("[SendMail] err: %v", err)
		return
	}

	mailID, _ = device.GenerateBizId(types.MailEntityBiz)
	mailEntity := models.MailEntity{
		Id:    mailID,
		Owner: senderID,
		Title: title,
		Body:  body,
		Ctime: libtools.GetUnixMillis(),
	}
	_, err = models.OrmInsert(&mailEntity)
	if err != nil {
		logs.Error("[SendMail] insert data get exception, mailEntity: %#v, err: %v", mailEntity, err)
		return
	}
	for _, receiverID := range receiverIDBox {
		mailBoxID, _ := device.GenerateBizId(types.MailboxBiz)
		mailbox := models.Mailbox{
			Id:         mailBoxID,
			MailID:     mailID,
			SenderID:   senderID,
			ReceiverID: receiverID,
			ReceiveAt:  libtools.GetUnixMillis(),
		}
		_, err = models.OrmInsert(&mailbox)
		if err != nil {
			// 有可能出现重复发送的,忽略就好
			logs.Warning("[SendMail] insert data get exception, mailbox: %#v, err: %v", mailbox, err)
			continue
		} else {
			mailboxList = append(mailboxList, mailBoxID)
		}
	}

	return
}

func List(condBox map[string]interface{}, page int, pageSize int) (list []ListItem, total int64, err error) {
	obj := models.Mailbox{}
	o := orm.NewOrm()

	var sql string
	var whereBox []string

	sqlCount := "SELECT COUNT(id) AS total"
	sqlQuery := fmt.Sprintf("SELECT %s.*, %s.title, %s.body",
		models.MAILBOX_TABLENAME, models.MAIL_ENTITY_TABLENAME, models.MAIL_ENTITY_TABLENAME)

	from := fmt.Sprintf("FROM %s", obj.TableName())
	joinFrom := fmt.Sprintf(`FROM %s LEFT JOIN %s ON %s.id = %s.mail_id`,
		models.MAILBOX_TABLENAME, models.MAIL_ENTITY_TABLENAME, models.MAIL_ENTITY_TABLENAME, models.MAILBOX_TABLENAME)

	var where string

	if v, ok := condBox["receiver_id"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`receiver_id = %d`, v.(int64)))
	}

	if _, ok := condBox["unread"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`read_at = 0`))
	}

	if v, ok := condBox["ctime_start_time"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`receive_at >= %d`, v.(int64)))
		whereBox = append(whereBox, fmt.Sprintf(`receive_at <= %d`, condBox["ctime_end_time"].(int64)))
	} else {
		whereBox = append(whereBox, fmt.Sprintf(`receive_at >= %d`, libtools.NaturalDay(-14)))
	}

	if _, ok := condBox["all_data"]; !ok {
		whereBox = append(whereBox, fmt.Sprintf(`%s.delete_at = 0`, models.MAILBOX_TABLENAME))
	}

	if len(whereBox) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereBox, " AND "))
	}

	var orderBy = "ORDER BY id DESC"
	if v, ok := condBox["order_by"]; ok {
		orderBy = v.(string)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = types.DefaultPagesize
	}
	offset := (page - 1) * pageSize

	limit := fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)

	sql = fmt.Sprintf("%s %s %s", sqlCount, from, where)
	_ = o.Raw(sql).QueryRow(&total)

	sql = fmt.Sprintf("%s %s %s %s %s", sqlQuery, joinFrom, where, orderBy, limit)
	logs.Notice("sql: %s", sql)
	_, err = o.Raw(sql).QueryRows(&list)

	return
}
