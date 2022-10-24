package types

//! 作务退出命令号
const TaskExitCmd int64 = -111

const XTrace = `111`

//! 金币系统的汇率 0.01元 == 1分
//! 数据库系统中所有钱相关的单位: 分

//!!! 非常重要: 结算的费率基数是 10000
// 1% -> 100, 1‰ -> 10, 1‱ -> 1
const FeeRateBase float64 = 10000

const CustomerServiceQr = "https://www.baidu.com"

const (
	DayMillisecond int64 = 86400000 // 1天的毫秒数
	DaySecond      int64 = 86400    // 1天的秒数

	OpenOauthWeChat   int = 1
	OpenOauthAlipay   int = 2
	OpenOauthFacebook int = 3
	OpenOauthGoogle   int = 4
	OpenOauthGuest    int = 5
)

func SuccessEmptyResp() struct{} {
	return struct{}{}
}

type BigInt int64

// Undefined 未定义
const Undefined = "未定义"

type PublicStatusEnum int

const (
	PublicStatusNo  PublicStatusEnum = 0
	PublicStatusYes PublicStatusEnum = 1
)

func PublicStatusEnumMap() map[PublicStatusEnum]string {
	return map[PublicStatusEnum]string{
		PublicStatusNo:  "公示中",
		PublicStatusYes: "已完成",
	}
}

type GenderEnum int

const (
	GenderUnknown GenderEnum = -1
	GenderFemale  GenderEnum = 0
	GenderMale    GenderEnum = 1
	GenderSecrecy GenderEnum = 2
)

func GenderEnumMap() map[GenderEnum]string {
	return map[GenderEnum]string{
		GenderUnknown: "未知",
		GenderFemale:  "女",
		GenderMale:    "男",
		GenderSecrecy: "保密",
	}
}

type MaritalStatusEnum int

const (
	MaritalStatusUnknown  MaritalStatusEnum = 0
	MaritalStatusSingle   MaritalStatusEnum = 1
	MaritalStatusMarried  MaritalStatusEnum = 2
	MaritalStatusDivorced MaritalStatusEnum = 3
	MaritalStatusWidowed  MaritalStatusEnum = 4
	MaritalStatusSecrecy  MaritalStatusEnum = 5
)

func MaritalStatusEnumConf() map[MaritalStatusEnum]string {
	return map[MaritalStatusEnum]string{
		MaritalStatusUnknown:  "未知",
		MaritalStatusSingle:   "单身",
		MaritalStatusMarried:  "已婚",
		MaritalStatusDivorced: "离异",
		MaritalStatusWidowed:  "丧偶",
		MaritalStatusSecrecy:  "保密",
	}
}

// DefaultPagesize 后台分页列表中的默认单页条数
const DefaultPagesize = 10

// admin session keys config
const (
	SessAdminIsLogin  string = "AdminIsLogin"
	SessAdminUid      string = "AdminUid"
	SessAdminNickname string = "AdminNickname"
	SessEditSaveGoto  string = `EditSaveGoto`
)

// orm 已经注册数据别名,需要有个`default`
const (
	OrmDataBase string = "default"
)

// 默认超管uid
const SuperAdminUID int64 = 1

type IdsBoxItem struct {
	Id int64
}

const (
	EmptyOrmStr       = "<QuerySeter> no row found"
	RedigoNilReturned = `redigo: nil returned`
)

// 平台定义
const (
	PlatformH5            = "h5"
	PlatformWxMiniProgram = "wx-mini-program"
	PlatformAndroid       = "android"
	PlatformAdm           = "web-adm"
	PlatformWatch         = `watch`
)

const (
	EventTaskRdsKey = "tinypro:queue:event"

	AdmLoginCaptchaCookieName = "_adm_login_captcha_"
)

type H map[string]interface{}

type MediaBaseInfo struct {
	MediaType   string
	Duration    int64
	DurationHum string
	Bitrate     int
	BitrateHum  string
	Width       int
	Height      int
	FirstFrame  []byte
}

type MediaSimpleInfo struct {
	MediaType string `json:"media_type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

type TagTwoTupleS struct {
	SN   string `json:"sn"`
	Name string `json:"name"`
}

type TwoTuple struct {
	SN   int64  `json:"sn,string"`
	Name string `json:"name"`
}

type OpMsgItem struct {
	OpBy    int64  `json:"op_by"`
	OpAt    int64  `json:"op_at"`
	Content string `json:"content"`
}

type WxPayTradeTypeEnum string
