package models

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
	"tinypro/common/types"
)

// 辅助方法, 为减少 orm 的初始化 和 m.Using 操作
// 减除冗余代码

// OrmModelPt ...
type OrmModelPt interface {
	TableName() string
}

func OrmOneByPkId(pkId int64, m OrmModelPt) (err error) {
	o := orm.NewOrm()

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE id = ? LIMIT 1`, m.TableName())
	err = o.Raw(sql, pkId).QueryRow(m)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[OrmOneByPkId] db fetch one exception, table: %s, id: %d, err: %v", m.TableName(), pkId, err)
		return
	}

	return nil
}

func buildListCond(condBox map[string]interface{}) (whereBox []string, whereArgs []interface{}) {
	if v, ok := condBox["id"]; ok {
		whereBox = append(whereBox, `id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["id_offset"]; ok {
		whereBox = append(whereBox, `id < ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["id_lt"]; ok {
		whereBox = append(whereBox, `id < ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["device_watch_id__gt"]; ok {
		whereBox = append(whereBox, `watch_device_id > ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["id_offset_lt"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`id < %d`, v.(int64)))
	}

	if v, ok := condBox["id_offset_gt"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`id > %d`, v.(int64)))
	}

	if v, ok := condBox["id_box"]; ok {
		if box, yes := v.([]interface{}); yes {
			if len(box) > 0 {
				whereBox = append(whereBox, fmt.Sprintf(`id IN (%s)`, libtools.SqlPlaceholderWithArray(len(box))))
				for _, obj := range box {
					whereArgs = append(whereArgs, obj)
				}
			}
		}
	}

	if v, ok := condBox["user_id"]; ok {
		whereBox = append(whereBox, `user_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["finish_flag"]; ok {
		whereBox = append(whereBox, `finish_flag = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["rewards"]; ok {
		whereBox = append(whereBox, `rewards = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["rewards__gt"]; ok {
		whereBox = append(whereBox, `rewards > ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["week_num"]; ok {
		whereBox = append(whereBox, `week_num = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["order_id"]; ok {
		whereBox = append(whereBox, `order_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["case_id"]; ok {
		whereBox = append(whereBox, `case_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["status"]; ok {
		whereBox = append(whereBox, `status = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["name"]; ok {
		whereBox = append(whereBox, `name LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%s%%`, v.(string)))
	}

	if v, ok := condBox["inviter_code"]; ok {
		whereBox = append(whereBox, `inviter_code = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["short_name"]; ok {
		whereBox = append(whereBox, `short_name LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%s%%`, v.(string)))
	}

	// 时间区间搜索
	if f, ok := condBox["ctime_start_time"]; ok {
		if field, has := condBox["_time_field_"]; has {
			whereBox = append(whereBox, fmt.Sprintf(`%s >= ?`, field.(string)))
			whereArgs = append(whereArgs, f)
			whereBox = append(whereBox, fmt.Sprintf(`%s < ?`, field.(string)))
			whereArgs = append(whereArgs, condBox["ctime_end_time"])
		}
	}

	if v, ok := condBox["created_by"]; ok {
		whereBox = append(whereBox, fmt.Sprintf(`created_by = %d`, v.(int64)))
	}

	if v, ok := condBox["mobile"]; ok {
		whereBox = append(whereBox, `mobile LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, v.(string)))
	}

	if v, ok := condBox["nickname"]; ok {
		whereBox = append(whereBox, `nickname LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, v.(string)))
	}

	if v, ok := condBox["student_id"]; ok {
		whereBox = append(whereBox, `student_id LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, v.(string)))
	}

	if v, ok := condBox["gender"]; ok {
		whereBox = append(whereBox, `gender= ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["owner"]; ok {
		whereBox = append(whereBox, `owner LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, v.(string)))
	}

	if v, ok := condBox["first_char"]; ok {
		whereBox = append(whereBox, `first_char = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["subject"]; ok {
		f := libtools.AddSlashes(v.(string))
		whereBox = append(whereBox, `subject LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, f))
	}

	if v, ok := condBox["model"]; ok {
		f := libtools.AddSlashes(v.(string))
		whereBox = append(whereBox, `model LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, f))
	}

	if v, ok := condBox["imei"]; ok {
		f := libtools.AddSlashes(v.(string))
		whereBox = append(whereBox, `imei LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%%%s%%`, f))
	}

	if v, ok := condBox["drive_permit_id"]; ok {
		whereBox = append(whereBox, `drive_permit_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["tags_name"]; ok {
		whereBox = append(whereBox, `name LIKE ?`)
		whereArgs = append(whereArgs, fmt.Sprintf(`%s%%`, v))
	}

	if v, ok := condBox["adminUID"]; ok {
		whereBox = append(whereBox, `admin_uid = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["e_account_number"]; ok {
		whereBox = append(whereBox, `e_account_number = ?`)
		whereArgs = append(whereArgs, v)
	}

	// 后台
	if f, ok := condBox["opTable"]; ok {
		whereBox = append(whereBox, `op_table = ?`)
		whereArgs = append(whereArgs, f)
	}

	if f, ok := condBox["opCode"]; ok {
		whereBox = append(whereBox, `op_code = ?`)
		whereArgs = append(whereArgs, f)
	}

	if f, ok := condBox["relatedId"]; ok {
		whereBox = append(whereBox, `related_id = ?`)
		whereArgs = append(whereArgs, f)
	}

	if f, ok := condBox["opUid"]; ok {
		whereBox = append(whereBox, `op_uid = ?`)
		whereArgs = append(whereArgs, f)
	}

	if v, ok := condBox["data_id"]; ok {
		whereBox = append(whereBox, `data_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["data_id_offset_gt"]; ok {
		whereBox = append(whereBox, `data_id > ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["data_id_offset_lt"]; ok {
		whereBox = append(whereBox, `data_id < ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["tag_id"]; ok {
		whereBox = append(whereBox, `tag_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["biz_sn"]; ok {
		whereBox = append(whereBox, `biz_sn = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["attr_type"]; ok {
		whereBox = append(whereBox, `attr_type = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["mark_num"]; ok {
		whereBox = append(whereBox, `mark_num = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["mark_pid"]; ok {
		whereBox = append(whereBox, `(mark_pid = ? OR  mark_num = ?)`)
		whereArgs = append(whereArgs, v, v)
	}

	if v, ok := condBox["lang_type"]; ok {
		whereBox = append(whereBox, `lang_type = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["classify"]; ok {
		whereBox = append(whereBox, `classify = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["class_id"]; ok {
		whereBox = append(whereBox, `class_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["student_id"]; ok {
		whereBox = append(whereBox, `student_id = ?`)
		whereArgs = append(whereArgs, v)
	}

	if v, ok := condBox["mark_num_box"]; ok {
		if box, yes := v.([]interface{}); yes {
			if len(box) > 0 {
				whereBox = append(whereBox, fmt.Sprintf(`mark_num IN (%s)`, libtools.SqlPlaceholderWithArray(len(box))))
				for _, obj := range box {
					whereArgs = append(whereArgs, obj)
				}
			}
		}
	}

	if v, ok := condBox["student_id_box"]; ok {
		if box, yes := v.([]interface{}); yes {
			if len(box) > 0 {
				whereBox = append(whereBox, fmt.Sprintf(`student_id IN (%s)`, libtools.SqlPlaceholderWithArray(len(box))))
				for _, obj := range box {
					whereArgs = append(whereArgs, obj)
				}
			}
		}
	}

	return
}

func OrmList(m OrmModelPt, condBox map[string]interface{}, page int, pageSize int, withCount bool, list interface{}) (total int64, err error) {
	o := orm.NewOrm()

	var sql string

	whereBox, whereArgs := buildListCond(condBox)

	sqlCount := "SELECT COUNT(id) AS total"
	sqlQuery := fmt.Sprintf("SELECT *")
	from := fmt.Sprintf("FROM %s", m.TableName())

	var where string

	if len(whereBox) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(whereBox, " AND "))
	}

	var orderBy = "ORDER BY id DESC"
	if v, ok := condBox["order_by"]; ok {
		orderBy = v.(string)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = types.DefaultPagesize
	}
	offset := (page - 1) * pageSize

	limit := fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)

	if withCount {
		sql = fmt.Sprintf("%s %s %s", sqlCount, from, where)
		errCount := o.Raw(sql, whereArgs).QueryRow(&total)
		if errCount != nil {
			logs.Error("[COUNT-SQL] db get exception, SQL: %s, args: %v err: %v", sql, whereArgs, errCount)
		}
	}

	sql = fmt.Sprintf("%s %s %s %s %s", sqlQuery, from, where, orderBy, limit)
	logs.Debug("[SELECT-SQL]: %s", sql)
	_, err = o.Raw(sql, whereArgs).QueryRows(list)
	if err != nil {
		logs.Error("[SELECT-SQL] db get exception, SQL: %s, args: %v, err: %v", sql, whereArgs, err)
	}

	return
}

// OrmInsert orm 插入
// m 必须为 model struct 的指针
// id 为自增主键的值
func OrmInsert(m OrmModelPt) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.Insert(m)
	return
}

// cols必须要唯一键束缚
// 例如[]string{"user_id","master_id"},表中user_id和master_id设置了唯一键
func OrmInsertOrUpdate(m OrmModelPt, cols []string) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.InsertOrUpdate(m, cols...)

	return
}

// OrmAllUpdate 全字段更新
// 若要使用全字段更新, m 必须为刚读出来的即时数据, 否则容易出现其他并发更新被覆盖
// m 必须为 model struct 的指针
// num 为 Affected Rows
func OrmAllUpdate(m OrmModelPt) (num int64, err error) {
	o := orm.NewOrm()

	num, err = o.Update(m)
	return
}

// OrmUpdate 特定字段更新
// 部分更新, 若是全字段更新, 请用 OrmAllUpdate
func OrmUpdate(m OrmModelPt, cols []string) (num int64, err error) {
	if len(cols) == 0 {
		// 部分更新, 必须指明字段, 防止错误更新
		err = fmt.Errorf("[OrmUpdate] can't do update with empty cols, %v", m)
		return
	}
	o := orm.NewOrm()

	num, err = o.Update(m, cols...)
	return
}

// OrmDelete 删除对象
func OrmDelete(m OrmModelPt) (num int64, err error) {
	o := orm.NewOrm()

	num, err = o.Delete(m)
	return
}

func OrmDeleteByCol(m OrmModelPt, cols []string) (num int64, err error) {
	o := orm.NewOrm()

	if len(cols) == 0 {
		num, err = o.Delete(m)
	} else {
		num, err = o.Delete(m, cols...)
	}
	return
}

func ParseBizSNFromPkID(pkID int64) (types.BizSN, error) {
	var err error
	bizI := pkID / 10000000000 % 100
	if bizI == 0 {
		err = fmt.Errorf(`input is 0`)
		logs.Warning("[ParseBizSNFromPkID] parse get exception, err: %v", err)
		return 0, err
	}

	biz := types.BizSN(bizI)
	return biz, err
}

func BizSN2TableName(bizSN types.BizSN) (string, error) {
	var name string
	var err error

	switch bizSN {

	case types.AccountSystem:
		name = ADMIN_TABLENAME

	case types.AppUserBiz:
		name = APP_USER_TABLENAME

	default:
		err = fmt.Errorf(`no register, bizSN: %d`, bizSN)
	}

	return name, err
}

func ParseTableNameFromPkID(pkID int64) (string, error) {
	var name string
	var err error

	bizSN, _ := ParseBizSNFromPkID(pkID)
	name, err = BizSN2TableName(bizSN)

	return name, err
}

func CheckDataValidByPkID(pkID int64) bool {
	tableName, err := ParseTableNameFromPkID(pkID)
	if err != nil {
		logs.Warning("[CheckDataValidByPkID] can parse table name, err: %v", err)
		return false
	}

	o := orm.NewOrm()

	var id int64
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE id = %d LIMIT 1`, tableName, pkID)
	err = o.Raw(sql).QueryRow(&id)
	if err != nil {
		logs.Error("[CheckDataValidByPkID] db get exception, SQL: %s, err: %v", sql, err)
		return false
	} else if id == pkID {
		return true
	}

	return false
}

func IncrByDistributedId(pkId int64, cols ...string) error {
	if pkId <= 0 || len(cols) == 0 {
		err := fmt.Errorf("pkID is 0, or cols is empty")
		logs.Warning("[IncrByPkID]", err)
		return err
	}

	table, err := ParseTableNameFromPkID(pkId)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	var setBox []string
	for _, f := range cols {
		setBox = append(setBox, fmt.Sprintf("%s = %s + 1", f, f))
	}

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d`, table, strings.Join(setBox, ", "), pkId)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[IncrByDistributedId] update get exception, sql: %s, err: %v", sql, err)
		return err
	}

	return nil
}

func IncrByPkId(table string, pkId int64, cols ...string) error {
	if pkId <= 0 || len(cols) == 0 {
		err := fmt.Errorf("pkID is 0, or cols is empty")
		logs.Warning("[IncrByPkID]", err)
		return err
	}

	o := orm.NewOrm()

	var setBox []string
	for _, f := range cols {
		setBox = append(setBox, fmt.Sprintf("%s = %s + 1", f, f))
	}

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d`, table, strings.Join(setBox, ", "), pkId)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[IncrByPkId] update get exception, sql: %s, err: %v", sql, err)
		return err
	}

	return nil
}
