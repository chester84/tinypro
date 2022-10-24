package subscribe_template_biz

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"tinypro/common/pkg/weixin/subscribe_template"
	"tinypro/common/pogo/reqs"
)

func SubscribeTemplate(userObj models.AppUser, req reqs.SubscribeTmpl) (err error) {
	cacheClient := storage.RedisStorageClient.Get()
	defer cacheClient.Close()

	for _, tmplId := range req.Ids {
		if subscribe_template.CheckTmplMapKey(tmplId) != "" {
			key := subscribe_template.GetRDSCourseTmplIDKey(tmplId, req.CourseSN)
			_, err = cacheClient.Do("SADD", key, userObj.Id)
			if err != nil {
				logs.Warning("SubscribeTemplate SADD err, cKey: %s, err %#v", key, err)
				return
			}
		} else {
			err = fmt.Errorf("tmplId not exist, tmplId: %s", tmplId)
			return
		}
	}
	return
}
