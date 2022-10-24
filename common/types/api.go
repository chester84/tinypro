package types

const (
	AppForceUpgradeConfigKey = `app_force_upgrade_config`
)

type ApiBaseT struct {
	TraceID     string `json:"trace_id"`
	RequestTime int64  `json:"request_time"`
	Data        string `json:"data"`
	EncryptKey  string `json:"encrypt_key"`
}

type WebApiBaseT struct {
	TraceID     string `json:"trace_id"`
	AccessToken string `json:"access_token"`
	RequestTime string `json:"request_time"`
	Imei        string `json:"imei"`
}

type ApiAppForceUpgradeT struct {
	OpMsg      string `json:"op_msg"`
	UpgradeMsg string `json:"upgrade_msg"` // 升级提示语
	ApkUrl     string `json:"apk_url"`     // apk下载链接
}

type ApiMobileOauthLoginReqT struct {
	Mobile string `json:"mobile"`
}

type ApiOauthLoginReqT struct {
	AppSN        int    `json:"app_sn"`
	OpenOauthPlt int    `json:"open_oauth_plt"`
	OpenUserID   string `json:"open_user_id"`
	Nickname     string `json:"nickname"`
	OpenAvatar   string `json:"open_avatar"`
	OpenToken    string `json:"open_token"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`

	WxOpenId string `json:"wx_open_id"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`

	Gender GenderEnum `json:"gender"`
}

type AppForceUpgradeConfItem struct {
	// 一期只做单个指定
	ConfType int `json:"version_type"` // 配置类型.1: 单个指定; 2: 区间

	AppVersion string `json:"app_version"` // 4位字符串版本 1.1234.1234.1234
	NumVersion int64  `json:"num_version"`

	UpgradeMsg string `json:"upgrade_msg"` // 升级提示语
	ApkUrl     string `json:"apk_url"`     // apk下载链接

	LastOpBy int64 `json:"last_op_by"`
	LastOpAt int64 `json:"last_op_at"`

	//BeginVersion string `json:"begin_version"`
	//BeginNum     int64  `json:"begin_num"`
	//EndVersion   string `json:"end_version"`
	//EndNum       int64  `json:"end_num"`
}

type ApiFaqItem struct {
	Subject string `json:"subject"` // 主题
	Content string `json:"content"` // 具休内容

	Tags []TagTwoTupleS `json:"tags"`
}

type ApiFaqCategoryItem struct {
	Category string       `json:"category"`
	List     []ApiFaqItem `json:"list"`
}
