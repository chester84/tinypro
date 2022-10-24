package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	// 数据库初始化
	_ "tinypro/common/lib/db/mysql"
	"tinypro/common/lib/redis/storage"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

var userID int64
var function string
var days int

func init() {
	//flag.Int64Var(&userID, "userID", 0, "传入正确的userid")
	//flag.IntVar(&days, "days", 0, "传入想要增加多少天的数据")
	//flag.StringVar(&function, "function", "", "传入正确的函数")
}

func main() {

}

func test() {
	logs.Debug("just one test")
}

func clearAllUserJoinData() {

	o := orm.NewOrm()

	//storageClient := storage.RedisStorageClient.Get()
	//defer storageClient.Close()

	r := o.Raw("TRUNCATE app_user")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE account_token")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_join_dare")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_upload_data")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_week_finish_task")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_mend_record")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_risk_manage")
	_, _ = r.Exec()

	r = o.Raw("TRUNCATE watch_risk_same_speed")
	_, _ = r.Exec()

	r = o.Raw("UPDATE watch_device SET back_live_card_send_time = 0, back_live_card_use_time = 0, stop_flag = 0")
	_, _ = r.Exec()

	//var userObj models.AppUser
	//var userList []models.AppUser
	//_, _ = o.QueryTable(userObj.TableName()).All(&userList)
	//
	//for _, oneUser := range userList {
	//
	//	rdsKey := join_dare.ChkInRdsKey(oneUser.Id)
	//	logs.Debug("rdsKey %s", rdsKey)
	//
	//	days := []int64{20210920, 20210921, 20210922, 20210923, 20210924, 20210925, 20210926, 20210927, 20210928, 20210929, 20210930}
	//
	//	for _, day := range days {
	//		_, err := storageClient.Do("SREM", rdsKey, day)
	//		if err != nil {
	//			logs.Error("SREM rdsKey %#v", err)
	//		}
	//
	//		humanDate := fmt.Sprintf("%d", day)
	//		rdsUserSuccessKey := join_dare.ChkInUserSuccessDayRdsKey(humanDate)
	//		logs.Debug("rdsUserSuccessKey %s", rdsUserSuccessKey)
	//
	//		_, err = storageClient.Do("SREM", rdsUserSuccessKey, oneUser.Id)
	//		if err != nil {
	//			err = fmt.Errorf("SREM rdsUserSuccessKey %#v", err)
	//			return
	//		}
	//	}
	//}

}

func clear20211005multi() {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	list, err := redis.Strings(storageClient.Do("SMEMBERS", "dn:checkin:users:success:20211005"))
	if err != nil {
		err = fmt.Errorf("SMEMBERS rdsUserSuccessKey %#v", err)
		return
	}

	for _, item := range list {
		hk := fmt.Sprintf("dn:hash:risk-rules:sp:%s", item)

		logs.Debug("[clear20211005multi] item %s", item)
		logs.Debug("[clear20211005multi] hk %s", hk)

		_, err = storageClient.Do("HDEL", hk, "2021-10-05")
		logs.Debug("[clear20211005multi] err %#v", err)
	}

}
