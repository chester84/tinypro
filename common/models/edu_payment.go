package models

import (
	"github.com/beego/beego/v2/client/orm"

	"tinypro/common/types"
)

func init() {
	orm.RegisterModel(new(EduPayment))
}

const EDU_PAYMENT_TABLENAME = "edu_payment"

type EduPayment struct {
	Id              int64 `orm:"pk;" json:"id,string"`
	AppSN           int   `orm:"column(app_sn)"` // 应用编号
	StudentId       int64
	ChargeConfId    int64
	FixedAmount     int64
	UnfixedAmount   int64
	Actual          int64
	Remark          string
	Status          types.PaymentStatusEnum
	CreatedAt       int64
	ClosedAt        int64
	WxPrepayId      string // 回填的微信支付凭据
	WxTransactionId string
	CallbackAt      int64
	LastOpBy        int64
	LastOpAt        int64
}

func (r *EduPayment) TableName() string {
	return EDU_PAYMENT_TABLENAME
}

type EduPaymentExpand struct {
	EduPayment

	ClassName     string
	CreatedUser   string
	StudentName   string
	Mobile        string
	Parent        string
	ParentContact string
	StudentSN     string
	ActualHuman   string
	LastOpUser    string
	StatusDesc    string
}
