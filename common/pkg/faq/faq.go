package faq

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"tinypro/common/types"
)

func ListInIds(ids []int64) (list []models.Faq, err error) {
	m := models.Faq{}
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableName()).
		Filter("id__in", ids).
		Filter("status", types.StatusValid).
		OrderBy("weight", "id").
		Limit(5).
		All(&list)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[ListInIds] db filter get exception, ids: %v, err: %v", ids, err)
	}

	return
}
