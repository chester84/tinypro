package models

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(new(EduPaymentExt))
}

const EDU_PAYMENT_EXT_TABLENAME = "edu_payment_ext"

type EduPaymentExt struct {
	Id        int64 `orm:"pk;" json:"id,string"`
	PaymentId int64
	ItemId    int64
	Amount    int64 // 实际金额
	Remark    string
	CreatedAt int64
}

func (r *EduPaymentExt) TableName() string {
	return EDU_PAYMENT_EXT_TABLENAME
}

type EduPaymentExtExpand struct {
	EduPaymentExt

	AmountHuman string
	ItemName    string
}
