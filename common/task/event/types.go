package event

import (
	"reflect"
)

// 事件有效期,毫秒
const (
	EventExpire   int64 = 30 * 60 * 1000
	EventExpireEX int64 = EventExpire / 1000
)

// QueueVal 描述触发事件在队列里的值
// 持久化参数结构体名, 会作为 EventName 存入队列, 方便获取
type QueueVal struct {
	EventName string `json:"n"`
	Data      []byte `json:"d"`
}

// GetStructName 获取 struct name , 以生成 map
func GetStructName(m interface{}) string {
	val := reflect.ValueOf(m)
	name := reflect.Indirect(val).Type().Name()
	//fmt.Println(name)
	return name
}

// -----

// demo演示事件
type DemoEv struct {
	DemoID int64 `json:"demo_id"`
	Time   int64 `json:"time"`
}
