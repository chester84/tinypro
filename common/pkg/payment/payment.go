package payment

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/device"
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func EduCreateOne(studentId int64, req types.EduPayCreateOrderReq) (eduPay models.EduPayment, err error) {
	var (
		checkTotal    int64
		fixedAmount   int64
		unfixedAmount int64
	)
	payId, _ := device.GenerateBizId(types.EduPaymentBiz)

	for _, item := range req.FixedGroup {
		var conf models.EduChargeItem
		cfgId, _ := libtools.Str2Int64(item.ItemSN)
		errC := models.OrmOneByPkId(cfgId, &conf)
		if errC != nil || conf.ChargeType != types.EduChargeTypeFixed {
			continue
		}

		checkTotal += conf.Amount
		fixedAmount += conf.Amount

		ext := models.EduPaymentExt{
			PaymentId: payId,
			ItemId:    conf.Id,
			Amount:    conf.Amount,
			Remark:    item.Remark,
			CreatedAt: libtools.GetUnixMillis(),
		}
		_, errI := models.OrmInsert(&ext)
		if errI != nil {
			logs.Error("[EduCreateOne] db insert fixed payment ext exception, req: %#v, ext: %#v, err: %v", req, ext, errI)
		}
	}

	for _, item := range req.UnfixGroup {
		var conf models.EduChargeItem
		cfgId, _ := libtools.Str2Int64(item.ItemSN)
		errC := models.OrmOneByPkId(cfgId, &conf)
		if errC != nil || conf.ChargeType != types.EduChargeTypeUnfixed {
			continue
		}

		checkTotal += item.AmountNum
		unfixedAmount += item.AmountNum

		ext := models.EduPaymentExt{
			PaymentId: payId,
			ItemId:    conf.Id,
			Amount:    item.AmountNum,
			Remark:    item.Remark,
			CreatedAt: libtools.GetUnixMillis(),
		}
		_, errI := models.OrmInsert(&ext)
		if errI != nil {
			logs.Error("[EduCreateOne] db insert unfixed payment ext exception, req: %#v, ext: %#v, err: %v", req, ext, errI)
		}
	}

	eduPay.Id = payId
	eduPay.StudentId = studentId
	eduPay.AppSN = req.AppSN
	eduPay.ChargeConfId, _ = libtools.Str2Int64(req.ChargeSN)
	eduPay.FixedAmount = fixedAmount
	eduPay.UnfixedAmount = unfixedAmount
	eduPay.Actual = req.TotalAmount
	eduPay.Remark = req.Remark
	eduPay.Status = types.PaymentStatusCreated
	eduPay.CreatedAt = libtools.GetUnixMillis()

	_, err = models.OrmInsert(&eduPay)
	if err != nil {
		logs.Error("[EduCreateOne] db insert get exception, one: %#v, err: %v", eduPay, err)
		return
	}

	if checkTotal != req.TotalAmount {
		err = fmt.Errorf(`金额结算不一致`)
		logs.Error("[EduCreateOne] check aont fail, stdId: %d, req: %#v", studentId, req)
		return
	}

	return
}

func EduWriteChargeItem(opUid, chargeConfId int64, chargeType types.EduChargeTypeEnum, snapBox []types.EduChargeItemSimple) {
	for _, snap := range snapBox {
		itemId, _ := libtools.Str2Int64(snap.ItemSN)
		if itemId <= 0 {
			continue
		}

		amount, _ := libtools.DecimalMoneyMul100(snap.Amount)
		var itemObj models.EduChargeItem
		errC := models.OrmOneByPkId(itemId, &itemObj)
		if errC != nil {
			// 没有查到,新增
			itemObj.Id = itemId
			itemObj.ChargeConfId = chargeConfId
			itemObj.Name = snap.Name
			itemObj.ChargeType = chargeType
			itemObj.Amount = amount
			itemObj.Status = types.StatusValid
			itemObj.CreatedBy = opUid
			itemObj.CreatedAt = libtools.GetUnixMillis()
			itemObj.LastOpBy = opUid
			itemObj.LastOpAt = libtools.GetUnixMillis()
			_, errI := models.OrmInsert(&itemObj)
			if errI != nil {
				logs.Error(`[EduWriteChargeItem] db insert charge item exception, data: %#v, err: %v`, itemObj, errI)
			}
		} else {
			// 之前有,更新
			itemObj.Name = snap.Name
			itemObj.Amount = amount
			itemObj.LastOpBy = opUid
			itemObj.LastOpAt = libtools.GetUnixMillis()
			_, errU := models.OrmUpdate(&itemObj, []string{`Name`, `Amount`, `LastOpBy`, `LastOpAt`})
			if errU != nil {
				logs.Error(`[EduWriteChargeItem] db update charge item exception, data: %#v, err: %v`, itemObj, errU)
			}
		}
	}
}
