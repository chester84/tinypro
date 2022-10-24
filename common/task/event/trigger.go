package event

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/lib/redis/storage"
	"tinypro/common/types"
)

// Trigger 触发事件
// persistentParam 必须是在 event/evtypes/persistent_param.go
// 异步运行方法必须定义在 runevent包中
// Event trigger, will be import by anywhere to trigger event
// calls events
func Trigger(persistentParam interface{}) (ok bool, err error) {
	if persistentParam == nil {
		err = fmt.Errorf("[event.Trigger] persistentParam can not be nil, persistentParam:%v", persistentParam)
		logs.Error(err)
		return
	}

	ok = false
	var eqv QueueVal
	eqv.EventName = GetStructName(persistentParam)
	eqv.Data, _ = json.Marshal(persistentParam)

	// 如果配置, 即时触发, 可以直接此处 run event

	storageClient := storage.RedisStorageClient.Get()
	defer storageClient.Close()

	data, _ := json.Marshal(eqv)

	key := types.EventTaskRdsKey
	_, err = storageClient.Do("LPUSH", key, data)
	if err != nil {
		logs.Error("[event.Trigger] Event Queue, LPUSH", err, "; Event: ", persistentParam)
		return
	}
	ok = true
	return
}
