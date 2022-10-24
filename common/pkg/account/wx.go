package account

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/models"
	"github.com/chester84/libtools"
)

func WriteWxOpenId(userId int64, appSN int, openId string) {
	o := orm.NewOrm()
	m := models.WxOpenId{}

	err := o.QueryTable(m.TableName()).
		Filter("user_id", userId).Filter("app_sn", appSN).Filter("open_id", openId).Limit(1).
		One(&m)
	if err != nil {
		if err != orm.ErrNoRows {
			logs.Error("[WriteWxOpenId] db get unknown exception, userID: %d, appSN: %d, openId: %s, err: %v",
				userId, appSN, openId, err)
		} else {
			// 新增
			m.UserId = userId
			m.AppSN = appSN
			m.OpenId = openId
			m.CreatedAt = libtools.GetUnixMillis()

			_, err = models.OrmInsert(&m)
			if err != nil {
				logs.Error("[WriteWxOpenId] db insert get exception, m: %#v, err: %v", m, err)
			}
		}
	} else {
		// 更新
		m.LastOpAt = libtools.GetUnixMillis()
		_, err := models.OrmUpdate(&m, []string{"LastOpAt"})
		if err != nil {
			logs.Error("[WriteWxOpenId] db update exception, m: %#v, err: %v", m, err)
		}
	}
}

func GetWxOpenId(userId int64, appSN int) (openId string, err error) {
	o := orm.NewOrm()
	m := models.WxOpenId{}

	err = o.QueryTable(m.TableName()).
		Filter("user_id", userId).Filter("app_sn", appSN).Limit(1).
		One(&m)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("[GetWxOpenId] filter data exception, userId: %d, appSN: %d, err: %v", userId, appSN)
	}

	openId = m.OpenId

	return
}
