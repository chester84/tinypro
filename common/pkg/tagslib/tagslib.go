package tagslib

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

func TagByName(name string) (one models.TagsLibrary, err error) {
	o := orm.NewOrm()

	err = o.QueryTable(one.TableName()).
		Filter("name", name).
		Filter("status", types.StatusValid).
		OrderBy("-id").
		Limit(1).
		One(&one)

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("[tagslib->TagByName] db filter get exception, name: %s, err: %v", name, err)
	}

	return
}

func DataTagListByTagBiz(tagId int64, bizSN types.BizSN) (list []models.DataTagMap, err error) {
	m := models.DataTagMap{}

	o := orm.NewOrm()

	_, err = o.QueryTable(m.TableName()).
		Filter("tag_id", tagId).
		Filter("biz_sn", bizSN).
		Filter("status", types.StatusValid).
		OrderBy("-id").
		All(&list)

	if err != nil {
		logs.Error("[tagslib->DataTagListByTagBiz] db filter get exception, tag_id: %d, err: %v", tagId, err)
	}

	return
}

// FilterDataIdsByTags 通过一组tagId过滤出所有符合条件的数据id
func FilterDataIdsByTags(tags []int64, bizSN types.BizSN) (idsBox []int64) {
	var box = make(map[int64]bool)

	for _, tagId := range tags {
		dataBox, err := DataTagListByTagBiz(tagId, bizSN)
		if err != nil {
			continue
		}

		for _, data := range dataBox {
			if !box[data.DataID] {
				box[data.DataID] = true
				idsBox = append(idsBox, data.DataID)
			}
		}
	}

	return
}

func Intersect(a, b []models.DataTagMap) (box []int64) {
	m := map[int64]bool{}

	for _, data := range a {
		m[data.DataID] = true
	}

	for _, data := range b {
		if m[data.DataID] {
			box = append(box, data.DataID)
		}
	}

	return
}

func TagTowTuple(id int64) (types.TagTwoTupleS, error) {
	var (
		tuple types.TagTwoTupleS
		one   models.TagsLibrary
		err   error
	)

	err = models.OrmOneByPkId(id, &one)
	if err != nil {
		return tuple, err
	}

	tuple.SN = libtools.Int642Str(one.Id)
	tuple.Name = one.Name

	return tuple, nil
}

func DataTagTupleGroup(dataID int64) []types.TagTwoTupleS {
	var g = make([]types.TagTwoTupleS, 0)

	var list []models.DataTagMap

	obj := models.DataTagMap{}
	o := orm.NewOrm()

	_, err := o.QueryTable(obj.TableName()).Filter("data_id", dataID).Filter("status", types.StatusValid).All(&list)
	if err != nil {
		logs.Error("[DataTagTupleGroup] db get exception, err: %v", err)
		return g
	}

	for _, item := range list {
		sub, err := TagTowTuple(item.TagID)
		if err != nil {
			continue
		}

		g = append(g, sub)
	}

	return g
}

func AllTagsWithStatus(status types.StatusCommonEnum) map[int64]string {
	all := map[int64]string{}

	condBox := map[string]interface{}{}
	if status > -1 {
		condBox["status"] = int(status)
	}

	var (
		obj  models.TagsLibrary
		list []models.TagsLibrary
	)
	_, _ = models.OrmList(&obj, condBox, 1, 1000, false, &list)
	for _, t := range list {
		all[t.Id] = t.Name
	}

	return all
}

func CreateDataTag(dataID int64, tagIDBox []int64, opBy int64, bizSN types.BizSN) {
	if dataID <= 0 {
		logs.Warning("[CreateDataTag] input data id is 0")
		return
	}

	for _, tagID := range tagIDBox {
		if tagID <= 0 {
			logs.Warning("[CreateDataTag] input tag id is 0")
			continue
		}
		one := models.DataTagMap{
			DataID:   dataID,
			TagID:    tagID,
			BizSN:    bizSN,
			Status:   1,
			LastOpBy: opBy,
			LastOpAt: libtools.GetUnixMillis(),
		}

		_, err := models.OrmInsert(&one)
		if err != nil {
			logs.Error("[CreateDataTag] inster get exception, one: %#v, err: %v", one, err)
		}
	}
}

func UpdateDataTag(dataID int64, tagIDBox []int64, opBy int64, bizSN types.BizSN) {
	if dataID <= 0 {
		logs.Warning("[UpdateDataTag] input data id is 0")
		return
	}

	o := orm.NewOrm()
	m := models.DataTagMap{}

	sql := fmt.Sprintf(`UPDATE %s
SET status = 0, last_op_by = %d, last_op_at = %d
WHERE data_id = %d`,
		m.TableName(), opBy, libtools.GetUnixMillis(), dataID)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		logs.Error("[UpdateDataTag] update get exception, SQL: %s, err: %v", sql, err)
	}

	for _, tagID := range tagIDBox {
		if tagID <= 0 {
			logs.Warning("[UpdateDataTag] input tag id is 0")
			continue
		}

		check := models.DataTagMap{}
		err := o.QueryTable(m.TableName()).Filter("data_id", dataID).Filter("tag_id", tagID).One(&check)
		if err != nil {
			if err.Error() == types.EmptyOrmStr { // 没有原始数据,新增
				CreateDataTag(dataID, []int64{tagID}, opBy, bizSN)
			} else {
				logs.Error("[UpdateDataTag] check process exception, err: %v", err)
			}
		} else {
			// 更新状态, TODO 操作历史
			check.Status = 1
			check.LastOpBy = opBy
			check.LastOpAt = libtools.GetUnixMillis()
			_, err = models.OrmUpdate(&check, []string{"status", "last_op_by", "last_op_at"})
			if err != nil {
				logs.Error("[UpdateDataTag] update get exception, check: %#v, err: %v", check, err)
			}
		}
	}
}
