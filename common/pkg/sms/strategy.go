package sms

import (
	"tinypro/common/lib/redis/cache"
	"github.com/chester84/libtools"
	"tinypro/common/types"

	"github.com/beego/beego/v2/core/logs"
)

// 限制策略
var limitStrategyServiceMap = map[types.SmsServiceTypeEnum]string{
	types.ServiceRegister: "service-register",
	types.ServiceLogin:    "service-login",
}

type StrategyType int

const (
	TimesMore        StrategyType = 1
	FrequencyTooHigh StrategyType = 2
	ExpireTime                    = 5 * 60 * 1000
)

var limitStrategyMobile = map[string]int64{
	"frequency": 10,       //! 临时改为 10,运营商最多只发10条
	"interval":  86400000, // 银行级别的,24小时内只能试6次
	"coolTime":  60000,    // 两次短信的间隔最短时间
}

var limitStrategyPassword = map[string]int64{
	"notify":   3,        // 密码输错3次后，提示错误信息
	"lock":     6,        // 密码输错6次后，锁定用户
	"interval": 86400000, // 用户锁定后，24小时内自动解锁
}

type LimitStrategy struct {
	Strategy    string
	Mobile      string
	ServiceType types.SmsServiceTypeEnum
	SmsType     types.SmsTypeEnum
}

// 返回为 true 说明中了限制策略
// 判断24小时内验证码次数
func commonStrategy(smsLimitStrategy LimitStrategy) bool {
	strategy := smsLimitStrategy.Strategy
	mobile := smsLimitStrategy.Mobile
	serviceType := smsLimitStrategy.ServiceType
	smsType := smsLimitStrategy.SmsType

	// 考虑以后加更多限制策略的...
	if strategy != "mobile" {
		return false
	}

	// 想扩展,就扩展它
	var cfg = limitStrategyMobile

	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cKey := buildLimitStrategyKey(serviceType, smsType, mobile)
	cValue, _ := cacheClient.Do("GET", cKey)
	if cValue == nil {
		cacheClient.Do("SET", cKey, "1", "PX", cfg["interval"])
	} else {
		// 好个坑
		value := string(cValue.([]byte))
		num, _ := libtools.Str2Int64(value)
		num++
		//cacheClient.Do("SET", cKey, libtools.Int642Str(num))
		cacheClient.Do("INCR", cKey)
		if num > cfg["frequency"] {
			logs.Warning("hit limit strategy, key:", cKey, ", value:", num)
			return true
		}
	}

	return false // 不使用限制策略
}

// 短信验证码频率检查
func smsFrequencyStrategy(smsLimitStrategy LimitStrategy) bool {
	strategy := smsLimitStrategy.Strategy
	mobile := smsLimitStrategy.Mobile
	serviceType := smsLimitStrategy.ServiceType
	//authCodeType := authCodeStrategy.AuthCodeTypeVal

	// 考虑以后加更多限制策略的...
	if strategy != "mobile" {
		return false
	}

	var cfg = limitStrategyMobile

	smsVerifyData, err := GetSmsByMobile(mobile, serviceType)
	curMill := libtools.GetUnixMillis()
	timeDiff := curMill - smsVerifyData.CreatedAt.UnixMilli()
	if err == nil && timeDiff <= cfg["coolTime"] {
		return true
	}

	return false
}

// (短信验证码/语音验证码)一天只能发6次，并且(短信验证码/语音验证码)发送间隔大于60秒
func MobileStrategy(mobile string, serviceType types.SmsServiceTypeEnum, smsType types.SmsTypeEnum) (smsStrategy StrategyType) {
	smsLimitStrategy := LimitStrategy{
		Strategy:    "mobile",
		Mobile:      mobile,
		ServiceType: serviceType,
		SmsType:     smsType,
	}

	if smsFrequencyStrategy(smsLimitStrategy) {
		smsStrategy = FrequencyTooHigh
		return
	}

	if commonStrategy(smsLimitStrategy) {
		smsStrategy = TimesMore
		return
	}

	return
}
