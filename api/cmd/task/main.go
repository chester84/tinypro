package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/task"
	"github.com/chester84/libtools"
	"tinypro/common/types"

	// 数据库初始化
	_ "tinypro/common/lib/db/mysql"
)

/**
跑批任务限定单机,单进程,多协程,使用redis队列做解耦合
1. 生成队列
2. 消费队列
3. 安全退出
*/

const (
	programName = "tinypro-task"
)

var taskName string
var help bool
var version bool

func init() {
	flag.StringVar(&taskName, "name", "", "crontab, cli or backend `task-name`, need assign.")
	flag.StringVar(&taskName, "n", "", "crontab, cli or backend `task-name`, need assign.")
	flag.BoolVar(&help, "h", false, "show usage and exit")
	flag.BoolVar(&version, "v", false, "show version and exit")

	// 改变默认的 Usage
	flag.Usage = usage
	dir, _ := config.String("log_dir")
	initLog(dir, "task")
}

func usage() {
	taskWork0Map := task.TaskWork0Map()
	var output string = fmt.Sprintf("%s version: %s/%s\n", programName, programName, types.TaskVersion)
	output = fmt.Sprintf("%sgit-head-hash: %s\n", output, libtools.GitRevParseHead())
	output = fmt.Sprintf("%s\nUsage: task [-hv] --name=TASK_NAME\n\n", output)

	var nameBox []string
	for name, _ := range taskWork0Map {
		nameBox = append(nameBox, name)
	}

	sort.Strings(nameBox)

	var need bool
	for _, name := range nameBox {
		var fixName string
		nameLen := len(name)
		if nameLen < 32 {
			fixName = fmt.Sprintf(`%s:%s`, name, strings.Repeat(" ", 32-nameLen))
		}
		output = fmt.Sprintf("%s  %s%s\n", output, fixName, taskWork0Map[name].Describe)
		if !need {
			need = true
		}
	}
	if need {
		output += "\n"
	}

	output = fmt.Sprintf("%sOptions:\n", output)

	_, _ = fmt.Fprintf(os.Stderr, output)
	flag.PrintDefaults()
	os.Exit(0)
}

func showVersion() {
	_, _ = fmt.Fprintf(os.Stderr, programName+` version: `+programName+`/`+types.TaskVersion+"\n")
	_, _ = fmt.Fprintf(os.Stderr, "git-head-hash: %s\n", libtools.GitRevParseHead())
	os.Exit(0)
}

func initLog(dir, name string) {
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
			"error", "warning",
		}
	} else {
		logsConfig["separate"] = []string{
			"emergency", "alert", "critical",
			"error", "warning", "notice",
			"info", "debug",
		}
	}
	logCfgBson, _ := json.Marshal(logsConfig)
	_ = logs.SetLogger(logs.AdapterConsole)
	_ = logs.SetLogger(logs.AdapterMultiFile, string(logCfgBson))

	logs.Debug("runmode: ", runmode, ", config:", string(logCfgBson))
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
	} else if version {
		showVersion()
	}

	// TODO: 如果需要通过命令行传参,此外的逻辑需要升级
	taskWork0Map := task.TaskWork0Map()
	// fmt.Print(taskWork0Map)
	// os.Exit(0)

	if _, ok := taskWork0Map[taskName]; !ok {
		usage()
	}

	// 派发任务
	taskWorker := taskWork0Map[taskName]

	f := func() {
		str := task.GetCurrentData()
		logs.Warn("[main] handling signal current data:%s", str)
		taskWorker.Task.Cancel()
	}

	libtools.ClearOnSignal(f)

	taskWorker.Task.Start()
}
