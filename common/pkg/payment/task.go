package payment

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func EduFetchNeedCloseOrder() (box []int64) {
	var (
		err  error
		list []models.EduPayment
	)

	o := orm.NewOrm()
	queryTime := libtools.GetUnixMillis() - libtools.MillsSecondAHour
	_, err = o.QueryTable(models.EDU_PAYMENT_TABLENAME).
		Filter(`status`, types.PaymentStatusCreated).
		Filter(`created_at__lt`, queryTime).
		Filter(`closed_at`, 0).Limit(100).All(&list)
	if err != nil {
		logs.Error(`[EduFetchNeedCloseOrder] db filter data exception, err: %v`, err)
	}

	for _, payObj := range list {
		box = append(box, payObj.Id)
	}

	return
}
