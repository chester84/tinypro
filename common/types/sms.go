package types

// 短信运营商
type SmsVendorEnum int

const (
	Tencent SmsVendorEnum = 1
)

func SmsVendorEnumConf() map[SmsVendorEnum]string {
	return map[SmsVendorEnum]string{
		Tencent: "tecent",
	}
}

// 业务类型
type SmsServiceTypeEnum int

const (
	ServiceRegister   SmsServiceTypeEnum = 1
	ServiceLogin      SmsServiceTypeEnum = 2
	ServiceBindMobile SmsServiceTypeEnum = 3
)

func ServiceTypeEnumConf() map[SmsServiceTypeEnum]string {
	return map[SmsServiceTypeEnum]string{
		ServiceRegister:   "register",
		ServiceLogin:      "login",
		ServiceBindMobile: "bind-mobile",
	}
}

// 短信类型
type SmsTypeEnum int

const (
	SmsText  SmsTypeEnum = 1
	SmsAudio SmsTypeEnum = 2
)

var smsVerifyMap = map[SmsTypeEnum]string{
	SmsText:  "text",
	SmsAudio: "audio",
}

func SmsTypeEnumConf() map[SmsTypeEnum]string {
	return smsVerifyMap
}

// 短信发送状态
type SmsSendStatusEnum int

const (
	SmsSendFail    SmsSendStatusEnum = 0
	SmsSendSuccess SmsSendStatusEnum = 1
)

func SendStatusEnumConf() map[SmsSendStatusEnum]string {
	return map[SmsSendStatusEnum]string{
		SmsSendFail:    "sendFail",
		SmsSendSuccess: "sendSuccess",
	}
}

// 短信发送状态
type SmsVerifyStatusEnum int

const (
	SmsVerifyStatusFail        SmsVerifyStatusEnum = 0
	SmsVerifyStatusEnumSuccess SmsVerifyStatusEnum = 1
)

func SmsVerifyStatusEnumConf() map[SmsVerifyStatusEnum]string {
	return map[SmsVerifyStatusEnum]string{
		SmsVerifyStatusFail:        "verifyFail",
		SmsVerifyStatusEnumSuccess: "verifySuccess",
	}
}

// 短信发送状态
type SmsCallbackDeliveryStatusEnum int

const (
	SmsCallbackDeliveryFail    SmsCallbackDeliveryStatusEnum = 0
	SmsCallbackDeliverySuccess SmsCallbackDeliveryStatusEnum = 1
)

func CallbackDeliveryStatusEnumConf() map[SmsCallbackDeliveryStatusEnum]string {
	return map[SmsCallbackDeliveryStatusEnum]string{
		SmsCallbackDeliveryFail:    "smsCallbackDeliveryFail",
		SmsCallbackDeliverySuccess: "smsCallbackDeliverySuccess",
	}
}

const (
	DynamicCodeSmsTplID = "717653"
)
