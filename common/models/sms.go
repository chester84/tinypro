package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(Sms))
}

const APP_SMS_TABLENAME = "sms"

type Sms struct {
	Id                     int64 `orm:"pk;"`
	Vendor                 types.SmsVendorEnum
	ServiceType            types.SmsServiceTypeEnum
	Mobile                 string
	Code                   string
	Content                string
	SmsType                types.SmsTypeEnum
	Expires                int64
	SendStatus             types.SmsSendStatusEnum
	ResponseId             string
	Response               string
	Ip                     string
	CallbackDeliveryStatus types.SmsCallbackDeliveryStatusEnum
	CallbackContent        string
	VerifyStatus           types.SmsVerifyStatusEnum
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (r *Sms) TableName() string {
	return APP_SMS_TABLENAME
}
