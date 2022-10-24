// 事件消费集合

package runner

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/task/event"
	"github.com/chester84/libtools"
)

// event run place
// will be call by event consumer goroutine

// calls any outside and events

// to place all event here
// it will call by any where
// and call any where

func demoEv(param interface{}) (success bool, err error) {
	timeNow := libtools.GetUnixMillis()
	var ev *event.DemoEv

	if e, ok := param.(*event.DemoEv); ok {
		if timeNow-e.Time > event.EventExpire {
			logs.Informational("[demoEv] event has expired, DemoID: %d, time: %d", e.DemoID, e.Time)
			return
		}

		ev = e
	} else {
		err = fmt.Errorf("[demoEv] did not get a *event.DemoEv persistent param: %T", param)
		logs.Error(err)
		return
	}

	logs.Notice("event process get data, demoID: %d", ev.DemoID)

	return
}
