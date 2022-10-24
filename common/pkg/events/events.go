package events

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
)

func IncrByPkID(pkID int64, cols ...string) error {
	if pkID <= 0 || len(cols) == 0 {
		err := fmt.Errorf("pkID is 0, or cols is empty")
		logs.Warning("[IncrByPkID]", err)
		return err
	}

	table, err := models.ParseTableNameFromPkID(pkID)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var setBox []string
	for _, f := range cols {
		setBox = append(setBox, fmt.Sprintf("%s = %s + 1", f, f))
	}

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d`, table, strings.Join(setBox, ", "), pkID)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[IncrByPkID] update get exception, sql: %s, err: %v", sql, err)
		return err
	}

	return nil
}
