package museum

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"tinypro/common/types"
)

func ListInIds(ids []int64) (list []models.Museum, err error) {
	m := models.Museum{}
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableName()).
		Filter("id__in", ids).
		Filter("status", types.StatusValid).
		OrderBy("-id").
		All(&list)

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("[museum->ListInIds] db filter get exception, ids: %v, err: %v", ids, err)
	}

	return
}
