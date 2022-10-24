package runner

import (
	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
)

// ParseAndRun 解析 eventBytes,并运行事件
// 将 eventBytes 解析成 event , 根据已注册事件,找到其对应的 异步运行方法, 并运行
func ParseAndRun(eventBytes []byte) (success bool, err error) {
	// 崩溃log记录
	defer func() {
		if x := recover(); x != nil {
			logs.Error("[Run] panic data:%s, err:%v", string(eventBytes), x)
			logs.Error(libtools.FullStack())
		}
	}()

	return globalParser.Run(eventBytes)
}
