package system

import (
	"fmt"
	"tinypro/common/lib/redis/storage"
	"tinypro/common/models"
	"github.com/chester84/libtools"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
)

const (
	dateLayout = `Y-m-d`
)

// PvIncr 只管访问数,不管理具体是哪个接口
func PvIncr() {
	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	date := libtools.UnixMsec2Date(libtools.GetUnixMillis(), dateLayout)
	rdsKey := BuildPvKey(date)
	_, err := storageClient.Do("INCR", rdsKey)
	if err != nil {
		logs.Error(`[PvIncr] redis> INCR %s, err: %v`, rdsKey, err)
	}
}

// UvRecord 集合去重做uv统计
func UvRecord(userId int64) {
	if userId <= 0 {
		return
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	date := libtools.UnixMsec2Date(libtools.GetUnixMillis(), dateLayout)
	rdsKey := BuildUvBoxKey(date)
	_, err := storageClient.Do("SADD", rdsKey, userId)
	if err != nil {
		logs.Error(`[UvRecord] redis> SADD %s %d, err: %v`, rdsKey, userId, err)
	}
}

func PvUv(date string) (pv, uv int64) {
	var err error

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	pvKey := BuildPvKey(date)
	pv, err = redis.Int64(storageClient.Do("GET", pvKey))
	if err != nil && err != redis.ErrNil {
		logs.Error(`[PvUv] redis> GET %s, err: %v`, pvKey, err)
	}

	uvKey := BuildUvBoxKey(date)
	uv, err = redis.Int64(storageClient.Do(`SCARD`, uvKey))
	if err != nil {

	}

	return
}

func WriteKeepRecord(userObj *models.AppUser) {
	if userObj.Id <= 0 {
		return
	}

	var timeNow = libtools.GetUnixMillis()
	registerDate := libtools.UnixMsec2Date(userObj.RegisterAt.UnixMilli(), `Y-m-d`)
	registerNatural := libtools.Date2UnixMsec(registerDate, `Y-m-d`) // 注册当天的0时
	timeDiff := timeNow - registerNatural
	if timeDiff <= 0 {
		logs.Warning(`[WriteKeepRecord] register data exception, user: %#v`, userObj)
		return
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	date := libtools.UnixMsec2Date(timeNow, dateLayout)
	days := timeDiff / libtools.MillsSecondADay

	var keepKey string
	if days == 7 {
		keepKey = Build7KeepKey(date)
	} else if days == 2 {
		keepKey = Build2KeepKey(date)
	}

	if keepKey == "" {
		return
	}

	_, err := storageClient.Do("SADD", keepKey, userObj.Id)
	if err != nil {
		logs.Error(`[WriteKeepRecord] redis> SADD %s %d, err: %v`, keepKey, userObj.Id)
	}
}

func QueryKeepRecord(date string) (keep2, keep7 float64) {
	keep2 = 0
	keep7 = 0

	startAt := libtools.Date2UnixMsec(date, dateLayout)
	if startAt <= 0 {
		return
	}

	endAt := startAt + libtools.MillsSecondADay

	obj := models.AppUser{}
	o := orm.NewOrm()

	total, err := o.QueryTable(obj.TableName()).
		Filter(`register_at__gte`, startAt).
		Filter(`register_at__lt`, endAt).
		Count()
	if err != nil {
		logs.Error(`[QueryKeepRecord] db filter exception, err: %v`, err)
		return
	}

	if total <= 0 {
		logs.Info(`[QueryKeepRecord] no user join system, date: %s`, date)
		return
	}

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	kee2Key := Build2KeepKey(date)
	kee2Num, err := redis.Int64(storageClient.Do(`SCARD`, kee2Key))
	if err != nil {
		logs.Error(`[QueryKeepRecord] redis> SCARD %s, err: %v`, kee2Key, err)
	}
	if kee2Num > 0 {
		keep2, _ = libtools.Str2Float64(fmt.Sprintf(`%.2f%%`, float64(kee2Num)*100/float64(total)))
	}

	kee7Key := Build7KeepKey(date)
	keep7Num, err := redis.Int64(storageClient.Do(`SCARD`, kee7Key))
	if err != nil {
		logs.Error(`[QueryKeepRecord] redis> SCARD %s, err: %v`, kee7Key, err)
	}
	if keep7Num > 0 {
		keep7, _ = libtools.Str2Float64(fmt.Sprintf(`%.2f%%`, float64(keep7Num)*100/float64(total)))
	}

	return
}
