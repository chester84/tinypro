package cerror

import (
	"github.com/chester84/libtools"
)

// ErrCode represents a specific error type in a error class.
// Same error code can be used in different error classes.
type ErrCode int

type ApiResponse struct {
	Code      ErrCode     `json:"code"`
	Message   string      `json:"message"`
	SeverTime int64       `json:"sever_time"`
	Data      interface{} `json:"data"`
}

type AjaxResponse struct {
	Code    ErrCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var EmptyData = struct{}{}

const (
	// CodeUnknown is for errors of unknown reason.
	CodeUnknown ErrCode = 0

	CodeSuccess ErrCode = 200

	// 500 服务不可用,存在严重问题的特殊码
	Stop4Upgrade     ErrCode = 500110 // 停机维护
	ServiceIsDown    ErrCode = 500111 // 服务完全不可用,后端服务发生严重问题,需要时间恢复,让客户提示用户稍后再试
	ServiceDbOpFail  ErrCode = 500112 // 数据库操作失败
	RedisCacheDown   ErrCode = 500113 // redis 缓存不可用
	RedisStorageDown ErrCode = 500114 // redis 存储不可用

	// 800800
	AppForceUpgrade ErrCode = 800800 // 客户端强制升级

	// 4xx 接口相关
	LostRequiredParameters ErrCode = 400110 // 缺少必要参数
	SignatureVerifyFail    ErrCode = 400111 // 验证签名失败
	LostAccessToken        ErrCode = 400112 // 缺少 token
	AccessTokenExpired     ErrCode = 400113 // token 过期
	InvalidAccessToken     ErrCode = 400114 // 无效的 token
	InvalidMobile          ErrCode = 400115 // 无效的手机号
	RequestApiTooMore      ErrCode = 400117 // 请求接口过于频繁
	ApiNotFound            ErrCode = 400118 // 接口不存在
	InvalidRequestData     ErrCode = 400119 // 无效的请求数据
	UserUnRegister         ErrCode = 400120 // 用户未注册

	LimitStrategyMobile ErrCode = 400121 // 手机号使用受到限制,请求验证码达到上限,目前是每种类型24小时内6次
	ServiceUnavailable  ErrCode = 400122 // 服务不可用
	FileTypeUnsupported ErrCode = 400123 // 文件类型不支持
	UploadResourceFail  ErrCode = 400124 // 上传资源操作失败
	PermissionDenied    ErrCode = 400125 // 操作系统文件权限不足
	InvalidAccount      ErrCode = 400126 // 无效账户
	MobileAlreadyBind   ErrCode = 400127 // 手机号已经被绑定/账户已经绑定过手机号
	UploadResourceLost  ErrCode = 400128 // 上传资源不存在
	NeedBindMobile      ErrCode = 400129 // 需要绑定手机号

	UploadRepeat          ErrCode = 400130 // 需要绑定手机号
	InvalidParameterValue ErrCode = 400131 // 无效的参数值
	RepeatRecharge        ErrCode = 400132 // 重复充值
	BalanceIsNotEnough    ErrCode = 400133

	SMSRequestMoreFrequency ErrCode = 400140 // 获取短信请求过于频繁
	InvalidPassword         ErrCode = 400141 // 密码无效
	InvalidOldPassword      ErrCode = 400142 // 旧密码无效
	AccountLocked           ErrCode = 400143 // 账号被锁定/封禁
	InvalidMobileOrPassword ErrCode = 400144 // 用户名或密码错误
	LimitIPWhitelist        ErrCode = 400145 // ip白名单限制
	RepeatedSubmitData      ErrCode = 400146 // 重复的提交数据

	IncompleteData ErrCode = 400150 // 资料不全

	RequestOutOfQuotas ErrCode = 400160 // 请求超出配额
	SmsSendFail        ErrCode = 400161 // 发送短信失败
	OrderNotExit       ErrCode = 400162 // 订单不存在
	NoMoreData         ErrCode = 400163 // 没有更多数据
	CanNotExecSearch   ErrCode = 400164 // 无法搜索

	// 项目 直接相关码
	RequestHasExpired        ErrCode = 400400 // 请求已过期
	RequestFrequently        ErrCode = 400401 // 请求过于频繁
	ParameterTypeWrong       ErrCode = 400402 // 参数类型有误
	AccessDenied             ErrCode = 400403 // 权限不足
	ParameterValueOutOfRange ErrCode = 400404 // 参数值不在文档指定的范围内
	SystemConfigurationErr   ErrCode = 400405 // 系统参数配置错误
	ClientTimeASync          ErrCode = 404406 // 客户端时间不对
	UnrecognizedDevice       ErrCode = 404407 // 无法识别的设备
	AccountBoundDevice       ErrCode = 404408 // 当前账户已经绑定过设备
	RewardCannotDivide       ErrCode = 404409 // 奖金池金额不合法
	MendCommonErr            ErrCode = 404410 // 修理手表通用err

	// 项目企业微信
	WorkWXRequireParaErr          ErrCode = 500008
	WorkWXRequireSnLost           ErrCode = 500009
	WorkWXDepartmentIdLost        ErrCode = 500100 // 企业微信部门id不合法
	WorkWXQrcodeNameLost          ErrCode = 500101 // 活码名称不合法
	WorkWXQrcodeDescLost          ErrCode = 500102 // 活码描述不合法
	WorkWXQrcodeTagLost           ErrCode = 500103 // 活码标签不合法
	WorkWXQrcodeContactLost       ErrCode = 500104 // 活码联系人不合法
	WorkWXQrcodeCreateErr         ErrCode = 500106 // 活码创建失败
	WorkWXQrcodeJsonErr           ErrCode = 500107 // 活码json解析失败
	WorkWXQrcodeAddContactErr     ErrCode = 500108 // 活码调用腾讯创建二维码接口失败
	WorkWXQrcodeDepartmentLost    ErrCode = 500109 // 活码部门名称不合法
	WorkWXQrcodeWXConfigIdLost    ErrCode = 500301 // 活码微信config_id不合法
	WorkWXQrcodeEditErr           ErrCode = 500302 // 活码调用腾讯编辑二维码接口失败
	WorkWXMessageSendMessErr      ErrCode = 500303 // 群发消息为空
	WorkWXMessageRichMessErr      ErrCode = 500304 // 群发消息富文本消息为空
	WorkWXMessagDepartmentJsonErr ErrCode = 500305 // 群发消息部门消息为空

	SmsCodeFailed ErrCode = 700005 // 验证码错误

	InvalidCmd    ErrCode = 600100
	InvalidData   ErrCode = 600101
	CallApiFail   ErrCode = 600102 // 调用接口发生异常
	InvalidConfig ErrCode = 600103

	// 三方接口出错
	GcpUnavailable      ErrCode = 700100
	TcVehicleLicenseOCR ErrCode = 700101
	WeixinOauth2Fail    ErrCode = 700102

	// 运动打卡错误
	JoinStopButHasBackLiveCard   ErrCode = 800100
	JoinStopButHasNoBackLiveCard ErrCode = 800101
	JoinOver2Times               ErrCode = 800102
	JoinTimesException           ErrCode = 800103
	JoinTimeNotValid             ErrCode = 800104
	JoinNotNeedBackLiveCard      ErrCode = 800105
	JoinForceStop                ErrCode = 800106
)

var errorMessageMap = map[ErrCode]string{
	CodeUnknown:              "unknown",
	CodeSuccess:              "success",
	ApiNotFound:              "api not found, please check out",
	Stop4Upgrade:             "service is stop for upgrade",
	ServiceUnavailable:       "back-end service is not available",
	ServiceIsDown:            "service is down",
	ServiceDbOpFail:          "db service exception",
	LostRequiredParameters:   "lost required parameters",
	InvalidRequestData:       "invalid request data",
	SignatureVerifyFail:      "signature verify fail",
	RequestHasExpired:        "request has expired",
	ParameterTypeWrong:       "parameter type has wrong",
	ParameterValueOutOfRange: "parameter value out of limit, please check it out",
	SystemConfigurationErr:   "system parameter configuration error",
	RequestApiTooMore:        "request interface is too frequent",
	AccessDenied:             "access denied",
	LimitIPWhitelist:         "limit ip whitelist",
	InvalidAccessToken:       "invalid access token",
	InvalidMobileOrPassword:  "invalid mobile or password",
	InvalidAccount:           `invalid account`,
	MobileAlreadyBind:        `mobile already bind`,
	AccountLocked:            "account locked",
	OrderNotExit:             "order not exit",
	NoMoreData:               `no more data`,
	CanNotExecSearch:         `can not exec search`,
	InvalidMobile:            "invalid mobile",
	SMSRequestMoreFrequency:  "sms frequency over limit",
	UploadRepeat:             "upload repeated",
	InvalidOldPassword:       "原始密码不正确",
	UserUnRegister:           "用户未注册",

	RequestOutOfQuotas: "request out of quotas",
	AppForceUpgrade:    `app need force upgrade`,

	InvalidCmd:         "invalid cmd",
	InvalidData:        "invalid data",
	AccessTokenExpired: "access token expired",
	CallApiFail:        "call api fail",
	InvalidConfig:      "invalid config parameters",
	IncompleteData:     "info incomplete",
	RepeatedSubmitData: `repeated submit data`,
	PermissionDenied:   `system permission denied`,

	GcpUnavailable:      "third party service not available",
	TcVehicleLicenseOCR: "third party service not available, VLOCR",
	WeixinOauth2Fail:    `weixin oauth2 fail`,

	UnrecognizedDevice: `unrecognized device`,
	AccountBoundDevice: `account bound device`,

	ClientTimeASync: "your time is invalid",

	SmsCodeFailed: "SMS verification has wrong",

	NeedBindMobile:     `need bind mobile`,
	RepeatRecharge:     `repeat recharge`,
	BalanceIsNotEnough: `balance is not enough`,

	WorkWXDepartmentIdLost: `部门id不合法`,

	WorkWXRequireParaErr:          `参数错误`,
	WorkWXRequireSnLost:           `活码sn缺失`,
	WorkWXQrcodeCreateErr:         `创建活码失败`,
	WorkWXQrcodeTagLost:           `活码标签不合法`,
	WorkWXQrcodeJsonErr:           `活码json解析失败`,
	WorkWXQrcodeAddContactErr:     `活码调用腾讯创建二维码接口失败`,
	WorkWXQrcodeWXConfigIdLost:    `活码微信config_id不合法`,
	WorkWXQrcodeEditErr:           `活码调用腾讯编辑二维码接口失败`,
	WorkWXMessageSendMessErr:      `群发消息为空`,
	WorkWXMessageRichMessErr:      `群发消息富消息为空`,
	WorkWXMessagDepartmentJsonErr: `群发消息部门为空`,

	JoinStopButHasBackLiveCard:   `昨天挑战失败，当前有复活卡`,
	JoinStopButHasNoBackLiveCard: `昨天挑战失败，并没有复活卡或者复活卡已经使用过`,
	JoinOver2Times:               `当天参加挑战次数不得超过两次`,
	JoinTimesException:           `先参加第2次，然后参加第一次，数据有问题`,
	JoinTimeNotValid:             `参加活动不在合法时间内`,
	JoinNotNeedBackLiveCard:      `用户参与应正常参与挑战，无需复活卡`,
	JoinForceStop:                `用户被强制下线`,
	RewardCannotDivide:           `奖金池金额不合法`,
}

func ErrorMessage(code ErrCode) string {
	if msg, ok := errorMessageMap[code]; ok {
		return msg
	} else {
		return "undefined"
	}
}

func BuildApiResponse(code ErrCode, data interface{}) ApiResponse {
	r := ApiResponse{
		Code:      code,
		Message:   ErrorMessage(code),
		SeverTime: libtools.GetUnixMillis(),
		Data:      data,
	}
	return r
}

func BuildAjaxResponse(code ErrCode, data interface{}) AjaxResponse {
	r := AjaxResponse{
		Code:    code,
		Message: ErrorMessage(code),
		Data:    data,
	}
	return r
}
