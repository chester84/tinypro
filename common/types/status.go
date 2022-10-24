package types

// 封禁等状态
type StatusBlockEnum int

const (
	Block   StatusBlockEnum = 0
	Unblock StatusBlockEnum = 1
)

func StatusBlockEnumConf() map[StatusBlockEnum]string {
	return map[StatusBlockEnum]string{
		Block:   "Block",  // 封禁
		Unblock: "Normal", // 正常
	}
}

// 状态的一些通用定义
type StatusCommonEnum int

const (
	StatusInvalid StatusCommonEnum = 0
	StatusValid   StatusCommonEnum = 1
	//StatusAudit   StatusCommonEnum = 2
)

func StatusCommonMap() map[StatusCommonEnum]string {
	return map[StatusCommonEnum]string{
		StatusInvalid: `无效`,
		StatusValid:   `有效`,
		//StatusAudit:   `审核中`,
	}
}

type PaymentStatusEnum int

const (
	PaymentStatusCreated PaymentStatusEnum = 1 // 已创建,待支付
	PaymentStatusSuccess PaymentStatusEnum = 2 // 支付成功
	PaymentStatusFailure PaymentStatusEnum = 3 // 支付失败
	PaymentStatusClosed  PaymentStatusEnum = 4 // 支付已关闭
)

func PaymentStatusConf() map[PaymentStatusEnum]string {
	return map[PaymentStatusEnum]string{
		PaymentStatusCreated: `待支付`,
		PaymentStatusSuccess: `支付成功`,
		PaymentStatusFailure: `支付失败`,
		PaymentStatusClosed:  `已关闭`,
	}
}

type WithdrawStatusEnum int

const (
	WithdrawStatusAudit      WithdrawStatusEnum = 1 // 审核中
	WithdrawStatusReject     WithdrawStatusEnum = 2 // 审核拒绝
	WithdrawStatusClosed     WithdrawStatusEnum = 3 // 已关闭
	WithdrawStatusInProgress WithdrawStatusEnum = 4 // 提现中
	WithdrawStatusFail       WithdrawStatusEnum = 5 // 提现失败
	WithdrawStatusSuccess    WithdrawStatusEnum = 6 // 提现成功
)

func WithdrawStatusConf() map[WithdrawStatusEnum]string {
	return map[WithdrawStatusEnum]string{
		WithdrawStatusAudit:      `审核中`,
		WithdrawStatusReject:     `审核拒绝`,
		WithdrawStatusClosed:     `已关闭`,
		WithdrawStatusInProgress: `提现中`,
		WithdrawStatusFail:       `提现失败`,
		WithdrawStatusSuccess:    `提现成功`,
	}
}
