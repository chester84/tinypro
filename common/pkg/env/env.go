// env 设置系统的环境变量

package env

import (
	"tinypro/common/pkg/system/config"
)

func GetEnvVariable(key, defValue string) (value string) {
	value = config.ValidItemString(key)
	if value == "" && defValue != "" {
		value = defValue
	}

	return
}
