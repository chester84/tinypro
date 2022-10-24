package oplog

import (
	"tinypro/common/models"
	"github.com/chester84/libtools"
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
)

// OpLoggerTableMap 定义了 op_logger 中包含的所有日志表名的列表
// 新增请直接附加
var OpLoggerTableMap = map[string]string{
	"orders": "orders",
}

// OpLoggerListStru 描述后台 op_logger 日志列表的结构
type OpLoggerListStru struct {
	Id        int64
	RelatedId int64
	OpUid     int64
	OpCode    models.OpCodeEnum
	OpTable   string
	Ctime     int64
}

func ConvertInt64tString(in string) string {
	var orgInt64Map = make(map[string]int64)
	json.Unmarshal([]byte(in), &orgInt64Map)
	logs.Info("orgInt64Map:%#v", orgInt64Map)

	var orgStrMap = make(map[string]interface{})
	json.Unmarshal([]byte(in), &orgStrMap)
	logs.Info("orgStrMap:%#v", orgStrMap)

	if len(orgInt64Map) == 0 {
		return in
	}

	for k, v := range orgInt64Map {
		if v > 0 {
			orgStrMap[k] = libtools.Int642Str(v)
		}
	}

	out, _ := json.Marshal(orgStrMap)
	return string(out)
}
