package clogs

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

func InitLog(dir, name string) {
	// 加入日志,先简单一点配置
	// 日志级别只在dev环境为Trace,其他环境均为Warning
	logs.EnableFuncCallDepth(true)
	var logsConfig = make(map[string]interface{})
	logname := dir + "/" + name + ".log"
	logsConfig["filename"] = logname
	logsConfig["rotate"] = false

	var runmode, _ = config.String("runmode")
	if "dev" != runmode {
		logsConfig["level"] = logs.LevelWarning
		logsConfig["separate"] = []string{
			"emergency", "alert", "critical",
			"error", "warning", "notice",
			"info", "debug",
		}
	} else {
		logsConfig["separate"] = []string{
			"emergency", "alert", "critical",
			"error", "warning", "notice",
			"info", "debug",
		}
	}
	logCfgBson, _ := json.Marshal(logsConfig)
	//_ = logs.SetLogger(logs.AdapterConsole)
	_ = logs.SetLogger(logs.AdapterMultiFile, string(logCfgBson))

	logs.Debug("runmode: ", runmode, ", config:", string(logCfgBson))
}
