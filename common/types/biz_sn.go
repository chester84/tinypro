package types

// 业务编号
type BizSN int

const (
	// 管理系统
	AccountSystem  BizSN = 60 // 后台管理帐户
	MailEntityBiz  BizSN = 61 // 邮件主体
	MailboxBiz     BizSN = 62 // 邮件箱
	TagsLibBiz     BizSN = 63 // 标签库
	AccessTokenBiz BizSN = 64

	UploadResource BizSN = 11 // 上传资源

	EduPaymentBiz BizSN = 24
	PaymentBiz    BizSN = 12

	// 用户
	AppUserBiz BizSN = 81
	Teacher    BizSN = 82
	Course     BizSN = 83
	VideoBase  BizSN = 84
	Video      BizSN = 85
	Price      BizSN = 86
)
