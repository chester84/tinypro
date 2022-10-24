package types

type WxInvokeTypeEnum int

const (
	RemindAttendClass WxInvokeTypeEnum = 1
	CheckPass         WxInvokeTypeEnum = 2
)

func WxInvokeTypeEnumMap() map[WxInvokeTypeEnum]string {
	return map[WxInvokeTypeEnum]string{
		RemindAttendClass: "上课提醒",
		CheckPass:         "审核通过提醒",
	}
}
