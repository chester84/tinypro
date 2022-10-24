package types

// 查博士第三方接口调用Enum
type WxCallbackReqTypeEnum int

const (
	WxCallbackUnifiedOrder WxCallbackReqTypeEnum = 1
	WxCallbackRefund       WxCallbackReqTypeEnum = 2
	WxMsgCallback          WxCallbackReqTypeEnum = 3
)
