package device

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/redis/storage"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

/**
生成形如: 180102981234567812,
0-5:  	YYMMDD
6-7:  	BizSN
8-15:	device seq id
16-17:  毫秒数最后2位
*/
func GenerateBizId(bizSN types.BizSN) (int64, error) {
	t := time.Now()
	nano := t.UnixNano()
	millis := nano / 1000000

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	idDeviceKey := IdDeviceHashKey

	field := libtools.GetDate(millis / 1000)
	id, err := storageClient.Do("HINCRBY", idDeviceKey, field, 1)
	if err != nil {
		logs.Error("[GenerateBizId] redis err %#v", err)
	}

	logs.Info("[GenerateBizId] bizSN: %d, millis: %d, field: %s,  id: %d", int(bizSN), millis, field, id)
	bizIdStr := fmt.Sprintf("%d%02d%02d%02d%08d%02d", t.Year()%100, t.Month(), t.Day(), bizSN, id.(int64)%100000000, millis%100)
	bizId, err := strconv.ParseInt(bizIdStr, 10, 64)

	return bizId, err
}
