package runner

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mohae/deepcopy"

	"tinypro/common/task/event"
)

// 主要分为 EventParser(事件解析器) 和 XxxxEv(一类事件)

// 触发某事件
// 初始化 e1 := XxxxEv{}
// Trigger(e1)

// 事件运行
// 任务自动从事件队列中取出事件, 由事件解析器, 解析并运行
// 从redis 获取 二进制 val []byte
// EventParser.Run(val)

// 此处 regEvent 作为一个注册事件结构定义
type regEvent struct {
	// RunFunc 为异步事件运行时调用的方法
	RunFunc func(persistentParam interface{}) (success bool, err error)
	// PersistentParam 为异步事件触发时定义并在运行时的持久化数据
	PersistentParam interface{} // must be put in events
}

// EventParser 全局事件解析器, 注册解析事件, 自动保存当前事件Map
var globalParser *parser

func init() {
	// 注册事件
	globalParser = new(parser)
	globalParser.Register(
		// Demo
		regEvent{demoEv, new(event.DemoEv)},
	)
}

// parser 描述解析器类型
type parser struct {
	// 保存所有注册事件映射
	//EventMap  map[string]EventInterface
	RegisteredEventMap map[string]regEvent
}

// Register 注册事件
func (p *parser) Register(dbrs ...regEvent) {
	if p.RegisteredEventMap == nil {
		p.RegisteredEventMap = make(map[string]regEvent)
	}
	for _, v := range dbrs {
		name := event.GetStructName(v.PersistentParam)
		if _, ok := p.RegisteredEventMap[name]; !ok {
			p.RegisteredEventMap[name] = v
		}
	}
}

// Run 解析器解析并运行事件
func (p *parser) Run(d []byte) (success bool, err error) {
	success = false

	logs.Debug("[event][runner] string queue value", string(d))

	var eql event.QueueVal
	_ = json.Unmarshal(d, &eql)
	if v, ok := p.RegisteredEventMap[eql.EventName]; ok {
		// 此处拷贝的指针类型
		realParam := deepcopy.Copy(v.PersistentParam)
		_ = json.Unmarshal(eql.Data, realParam)
		// 此处传递的也是指针类型的 struct
		success, err = v.RunFunc(realParam)
		if err != nil {
			logs.Informational("[event.Parser.Run]", err)
		}
	} else {
		logs.Error("[event.Run] unregistered event:", string(d))
	}
	return
}
