package weixin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chester84/libtools"

	"github.com/beego/beego/v2/core/logs"
)

func Signature(params map[string]interface{}, secret string) string {
	delete(params, "sign")
	paramLen := len(params)
	if paramLen <= 0 {
		logs.Warning("[wxnotify->Signature] params len is 0, params: %v", params)
		return ""
	}

	cntr := make([]string, paramLen)
	var i int = 0
	for k, _ := range params {
		cntr[i] = k
		i++
	}

	// 按字典序列排序
	sort.Strings(cntr)

	str := "" // 待签名字符串
	for i = 0; i < paramLen; i++ {
		key := cntr[i]
		value := libtools.Stringify(params[key])
		if len(value) > 0 {
			str += fmt.Sprintf("%s=%s&", key, value)
		}
	}
	str += fmt.Sprintf("key=%s", secret)
	logs.Debug("[Signature] need signature str:", str)

	str = strings.ToUpper(libtools.Md5(str))
	logs.Debug("[Signature] sign: %s", str)

	//mac := hmac.New(sha256.New, []byte(secret))
	//mac.Write([]byte(str))

	return str
}
