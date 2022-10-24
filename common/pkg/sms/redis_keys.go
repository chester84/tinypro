package sms

import (
	"fmt"
	"tinypro/common/types"
)

const (
	smsVerifyRdsPrefix = `tinypro:cache:sms-verify:`
	smsLimitStrategy   = `tinypro:cache:limit-strategy:`
)

func buildSmsVerifyKey(serviceType types.SmsServiceTypeEnum, mobile string) string {
	serviceTypeDesc := types.ServiceTypeEnumConf()[serviceType]
	return fmt.Sprintf("%s%s:%s", smsVerifyRdsPrefix, serviceTypeDesc, mobile)
}

func buildLimitStrategyKey(serviceType types.SmsServiceTypeEnum, smsType types.SmsTypeEnum, mobile string) string {
	serviceTypeDesc := types.ServiceTypeEnumConf()[serviceType]
	smsTypeDesc := types.SmsTypeEnumConf()[smsType]

	return fmt.Sprintf("%s%s:%s:%s", smsLimitStrategy, serviceTypeDesc, smsTypeDesc, mobile)
}
