package tools

import (
	"math"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/types"
)

const MinWithdrawAmount int64 = 600

func BuildWithdrawBody(total int64) (body types.WithdrawBody) {
	body.Total = total
	if total <= MinWithdrawAmount {
		logs.Warning("[BuildWithdrawBody] input out of  range, total: %d", total)
		return
	}

	fee := int64(math.Ceil(float64(total) * 6 / 1000)) // 千6,或6元
	if fee < MinWithdrawAmount {
		fee = MinWithdrawAmount
	}

	body.Amount = total - fee
	body.Fee = fee

	return
}

func BuildServiceFeeBody(amount int64) (body types.ServiceFeeBody) {
	body.Amount = amount

	// 扣除平台服务费 = max(互助金*10%, 100) + max(互助金*6‰, 6)
	//body.Fee = int64(math.Max(10000, math.Ceil(float64(amount)*10/100))) +
	//	int64(math.Max(float64(MinWithdrawAmount), float64(amount)*6/1000))

	body.PltFee = int64(math.Max(10000, math.Ceil(float64(amount)*10/100)))
	body.ChanFee = int64(math.Max(float64(MinWithdrawAmount), math.Ceil(float64(amount)*6/1000)))
	body.Fee = body.PltFee + body.ChanFee

	body.Actual = body.Amount - body.Fee
	if body.Actual < 0 {
		body.Actual = 0
	}

	return
}
