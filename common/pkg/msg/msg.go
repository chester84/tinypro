package msg

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/adapter/orm"
	"tinypro/common/models"
	"tinypro/common/pogo/reqs"
	"tinypro/common/pogo/resps"
	"github.com/chester84/libtools"
	"tinypro/common/types"
	"sort"
)

func MyMsgList(userObj models.AppUser, req reqs.PageInfo) (retList []resps.MyMsg, err error) {
	m := models.Msg{}
	o := orm.NewOrm()
	list := make([]models.Msg, 0)
	retList = make([]resps.MyMsg, 0)

	pageSize := req.Size
	if req.Size > 100 {
		pageSize = types.DefaultPagesize
	}

	userIdArr := []int64{userObj.Id, 0}

	if req.Type == 0 {
		if req.SN <= 0 {
			//默认取最新记录
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id__in", userIdArr).
				OrderBy("-id").
				Limit(pageSize).
				All(&list)
		} else {
			_, err = o.QueryTable(m.TableName()).
				Filter("user_id__in", userIdArr).
				Filter("id__gt", req.SN).
				OrderBy("id").
				Limit(pageSize).
				All(&list)
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].ID > list[j].ID
		})
	} else {
		_, err = o.QueryTable(m.TableName()).
			Filter("user_id__in", userIdArr).
			Filter("id__lt", req.SN).
			OrderBy("-id").
			All(&list)
	}

	if err != nil && err.Error() != types.EmptyOrmStr {
		logs.Error("MyMsgList db filter get exception, err: %v", err)
		return
	}

	for _, item := range list {
		data := resps.MyMsg{}
		data.SN = item.ID
		data.MsgType = item.MsgType
		data.MsgTypeDesc = types.MsgTypeMap()[item.MsgType]
		data.MsgContent = item.MsgContent
		data.CreatedAt = libtools.UnixMsec2Date(item.CreatedAt.UnixMilli(), "Y-m-d H:i")

		retList = append(retList, data)
	}

	return
}
