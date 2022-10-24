package types

type WithdrawBody struct {
	Total  int64 // 总金额
	Amount int64 // 可提现金额
	Fee    int64 // 手续费
}

type ServiceFeeBody struct {
	Amount int64 // 总金额
	Fee    int64 // 总的服务费
	Actual int64 // 实际金额

	PltFee  int64 // 平台服务费
	ChanFee int64 // 通道手续费
}

type TransSettleEnum int64

const (
	TransSettleOut TransSettleEnum = -1 // 出账
	TransSettleIn  TransSettleEnum = 1  // 入账
)

func TransSettleConf() map[TransSettleEnum]string {
	return map[TransSettleEnum]string{
		TransSettleOut: `出账`,
		TransSettleIn:  `入账`,
	}
}

type TransFlowEnum int

const (
	TransFlowRecharge       TransFlowEnum = 1
	TransFlowPlatServiceFee TransFlowEnum = 2
	TransFlowWithdraw       TransFlowEnum = 3
	TransFlowFee            TransFlowEnum = 4
	TransPutBack            TransFlowEnum = 5
	TransUserPlatServiceFee TransFlowEnum = 6
	TransTransferOut        TransFlowEnum = 7
	TransChanFee            TransFlowEnum = 8
	TransAdjustment         TransFlowEnum = 9
)

func TransFlowEnumConf() map[TransFlowEnum]string {
	return map[TransFlowEnum]string{
		TransFlowRecharge:       `用户充值`,
		TransFlowPlatServiceFee: `平台服务费`, // 以平台维度
		TransFlowWithdraw:       `提现`,
		TransFlowFee:            `手续费`,
		TransPutBack:            `结算充正`,
		TransUserPlatServiceFee: `用户服务费`, // 以用户维度
		TransTransferOut:        `中转出账`,
		TransChanFee:            `通道费`,
		TransAdjustment:         `结算平账`,
	}
}
