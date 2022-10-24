package teacher

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/adapter/orm"
	"tinypro/common/models"
	"tinypro/common/types"
)

func GetTeachersByIdS(ids []int64) (retList []models.Teacher, err error) {
	m := models.Teacher{}
	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableName()).
		Filter("id__in", ids).
		OrderBy("-id").
		All(&retList)

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("[course->FrontPage] db filter get exception, err: %v", err)
		return
	}

	return
}
