package stat

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func newCustomersNum(start, end int64) int64 {
	var (
		num int64
		err error
	)

	o := orm.NewOrm()
	m := models.AppUser{}

	num, err = o.QueryTable(m.TableName()).Filter("register_at__gte", start).Filter("register_at__lt", end).Count()
	if err != nil {
		logs.Error("[newCustomersNum] db filter exception, startAt: %s, endAt: %s err: %v",
			libtools.UnixMsec2Date(start, "Y-m-d"), libtools.UnixMsec2Date(end, "Y-m-d"), err)
	}

	return num
}

func newCustomersStatWithDay(startDay int64) types.HighChartsSpLine {
	var hc = types.HighChartsSpLine{}

	var i int64 = 1
	var series = types.HighChartsSeries{
		Name: "新增数",
	}

	start := libtools.NaturalDay(startDay)

	for ; i <= 7; i++ {
		hc.XAxis = append(hc.XAxis, libtools.UnixMsec2Date(start, "Y-m-d"))
		end := libtools.NaturalDay(startDay + i)

		total := newCustomersNum(start, end)
		series.Data = append(series.Data, total)

		start = end
	}

	hc.Series = append(hc.Series, series)

	return hc
}

func NewCustomersStatLastWeek() types.HighChartsSpLine {
	return newCustomersStatWithDay(-13)
}

func NewCustomersStat() types.HighChartsSpLine {
	return newCustomersStatWithDay(-6)
}
