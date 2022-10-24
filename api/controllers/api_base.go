package controllers

import (
	"encoding/json"

	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus"

	"tinypro/common/cerror"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/i18n"
	"tinypro/common/pkg/metrics"
	"tinypro/common/pkg/security"
	"tinypro/common/pkg/system/config"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

// APIBaseController 接口基类
type APIBaseController struct {
	beego.Controller

	AesCBCKey string
	AesCBCIV  string

	IP            string
	AppVersion    string
	CurrentRouter string // 当前路由

	XTrace      bool
	TraceID     string
	RequestTime int64
	Lang        string
	UILang      string

	RequestBody []byte
	RequestData string // 原始请求中`data`的明文json串

	AccountID int64 // app_user->id

	// request json,用于必要参数快捷检查,实际取参数,建议采用结构体
	RequestJSON map[string]interface{}
	beginTime   time.Time
}

func (c *APIBaseController) Prepare() {
	c.IP = c.Ctx.Input.IP()
	c.CurrentRouter = c.Ctx.Input.URL()
	c.Lang = types.LangEnUS
	c.UILang = types.LangEnUS

	// 量化接口并发量
	metrics.WebRequestTotal.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Inc()

	c.beginTime = time.Now()

	// 停机维护公告
	var stop4Upgrade bool
	if stop4Upgrade {
		c.CommonResp(cerror.Stop4Upgrade, cerror.EmptyData)
		return
	}

	//logs.Debug("ReqHeader:", c.Ctx.Request.Header)
	//// 需要配置文件配合 copyrequestbody = true

	userAgent := c.Ctx.Input.UserAgent()
	if userAgent == "" || !libtools.VerifyHttpUserAgent(userAgent) {
		logs.Warning("access denied, router: %s, UA: %s, ip: %s", c.CurrentRouter, userAgent, c.IP)
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	xTrace := c.Ctx.Request.Header.Get("X-Trace")
	if xTrace == types.XTrace {
		c.XTrace = true
	}

	c.RequestBody = c.Ctx.Input.RequestBody
	if !libtools.IsProductEnv() || c.XTrace {
		// 联调打印原始数据
		logs.Notice(">>> router: %s, ip: %s, RequestBody: %s", c.CurrentRouter, c.IP, string(c.RequestBody))
	}

	if len(c.RequestBody) < 16 {
		logs.Warning("post data is empty, router: %s, ip: %s", c.CurrentRouter, c.IP)
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var reqData types.ApiBaseT
	err := json.Unmarshal(c.RequestBody, &reqData)
	if err != nil {
		logs.Warning("cat NOT json decode request body, router: %s, ip: %s, data: %s", c.CurrentRouter, c.IP, string(c.RequestBody))
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	// 请求时间和服务器时间相差不能超过指定时间
	timeNow := libtools.GetUnixMillis()
	timeDiff := libtools.AbsInt64(timeNow - reqData.RequestTime)
	if reqData.TraceID == "" || reqData.EncryptKey == "" ||
		timeDiff > 300000 {
		logs.Warning("capture suspected attacks, router: %s, traceID: %s, ip: %s, requestTime: %d, sysTime: %d, timeDiff: %d, httpBody: %s",
			c.CurrentRouter, c.TraceID, c.IP, reqData.RequestTime, timeNow, timeDiff, string(c.RequestBody))
		c.TerminateWithCode(cerror.ClientTimeASync)
		return
	}

	c.TraceID = reqData.TraceID
	c.RequestTime = reqData.RequestTime
	logs.Info("[Prepare] client start request, router: %s, traceID: %s, ip: %s,  at: %s",
		c.CurrentRouter, c.TraceID, c.IP, libtools.UnixMsec2Date(c.RequestTime, "Y-m-d H:i:s"))

	// 抓包数据防重入
	if !security.PassStrongPreventRepeatedEntry(c.CurrentRouter, c.TraceID, c.IP) {
		logs.Warning("repeated entry, router: %s, traceID: %s, ip: %s, httpBody: %s",
			c.CurrentRouter, c.TraceID, c.IP, string(c.RequestBody))
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	encryptKey := reqData.EncryptKey
	rawData := reqData.Data

	aesInfo, errD := libtools.RsaDecrypt(encryptKey)
	if errD != nil {
		logs.Error("[RsaDecrypt] rsa decrypt exception, router: %s, traceID: %s, ip: %s, encryptKey: %s, err: %v",
			c.CurrentRouter, c.TraceID, c.IP, encryptKey, err)
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	aesArray := strings.Split(aesInfo, "|")
	if len(aesArray) != 2 {
		logs.Error("[Prepare] aes key is unexpected, router: %s, traceID: %s, ip: %s, encryptKey: %s, aesKey: %s",
			c.CurrentRouter, c.TraceID, c.IP, encryptKey, aesInfo)
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	c.AesCBCKey = aesArray[0]
	c.AesCBCIV = aesArray[1]
	if c.XTrace {
		logs.Notice("[trace] parse aes Key/IV, router: %s, traceID: %s, ip: %s, aesKey: %s, aesIV: %s",
			c.CurrentRouter, c.TraceID, c.IP, c.AesCBCKey, c.AesCBCIV)
	}

	// 解密data数据
	decryptData, desErr := libtools.AesDecryptCBC(rawData, aesArray[0], aesArray[1])
	if desErr != nil {
		logs.Warning("post data can NOT decrypt, router: %s, ip: %s, traceID: %s, data: %s, key: %s, iv: %s, err: %v",
			c.CurrentRouter, c.TraceID, c.IP, rawData, aesArray[0], aesArray[1], desErr)
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	c.RequestData = decryptData

	var reqJSON map[string]interface{}
	err = json.Unmarshal([]byte(decryptData), &reqJSON)
	if err != nil {
		logs.Warning("cat NOT json decode request decryptData, router: %s, traceID: %s, ip: %s, data:",
			c.CurrentRouter, c.TraceID, c.IP, decryptData)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	// json decode 通过
	c.RequestJSON = reqJSON

	// 必要参数检查,只检查存在,没有判值
	requiredParameter := map[string]bool{
		"noise":        true,
		"access_token": true,
		"app_sn":       true,
		"app_version":  true,
		"timezone":     true,
		"brand":        true,
		"network":      true,
		"market":       true,
		"lang":         true,
	}
	//! 通用的必传参数在此外做统一校验
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		logs.Warning("lost required parameters, router: %s, traceID: %s, ip: %s, request: %s",
			c.CurrentRouter, c.TraceID, c.IP, c.RequestData)
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	if lang, ok := c.RequestJSON["lang"].(string); ok {
		if types.Short2LangType(lang) != types.LanguageTypeNone {
			c.Lang = lang
		}
	}

	if uiLang, ok := c.RequestJSON["ui_lang"].(string); ok {
		if types.Short2LangType(uiLang) != types.LanguageTypeNone {
			c.UILang = uiLang
		}
	}

	// 以下是强制更新逻辑
	appVersion, _ := c.RequestJSON["app_version"].(string)
	c.AppVersion = appVersion
	appNumVer := libtools.AppNumVersion(appVersion)
	appForceUpgrade := types.ApiAppForceUpgradeT{
		OpMsg: `please update to latest version.`,
	}
	// 如果客户端传的版本号不符合文档要求,直接不让用
	if appNumVer <= 0 {
		c.CommonResp(cerror.AppForceUpgrade, appForceUpgrade)
		return
	}

	var appForceUpgradeConf []types.AppForceUpgradeConfItem
	confData := config.ValidItemString(types.AppForceUpgradeConfigKey)
	err = json.Unmarshal([]byte(confData), &appForceUpgradeConf)
	if err != nil {
		logs.Info("json decode get exception, router: %s, traceID: %s, ip: %s, confData: %s, err: %v",
			c.CurrentRouter, c.TraceID, c.IP, confData, err)
	}
	//logs.Debug("appNumVer: %d, appForceUpgradeConf: %#v", appNumVer, appForceUpgradeConf)
	for _, confItem := range appForceUpgradeConf {
		if confItem.NumVersion == appNumVer {
			// 此版本需要强制更新
			appForceUpgrade.OpMsg = `for better experience, please update to latest version`
			appForceUpgrade.UpgradeMsg = confItem.UpgradeMsg
			appForceUpgrade.ApkUrl = confItem.ApkUrl

			c.CommonResp(cerror.AppForceUpgrade, appForceUpgrade)
			return
		}
	}

	uri := c.Ctx.Request.RequestURI
	// 以下路由不需要持有 token
	notNeedTokenRoute := map[string]bool{
		"/api/v1/encrypt/ping": true,
	}
	if !notNeedTokenRoute[uri] {
		// 检查 token 有效性
		accessToken, ok := c.RequestJSON["access_token"].(string)
		if !ok {
			logs.Warning("access_token out of docs range, router: %s, traceID: %s, ip: %s, requestData: %s",
				c.CurrentRouter, c.TraceID, c.IP, c.RequestData)
			c.TerminateWithCode(cerror.InvalidAccessToken)
			return
		}

		ok, accountId := accesstoken.IsValidAccessToken(types.PlatformAndroid, accessToken)
		if !ok {
			logs.Warning("access token is invalid, router: %s, traceID: %s, ip: %s, requestData: %s",
				c.CurrentRouter, c.TraceID, c.IP, c.RequestData)
			c.TerminateWithCode(cerror.InvalidAccessToken)
			return
		}

		c.AccountID = accountId
	}

	if c.XTrace {
		logs.Notice("[xTrace] Prepare, router: %s, traceID: %s, ip: %s, request: %s",
			c.CurrentRouter, c.TraceID, c.IP, c.RequestData)
	}
}

func (c *APIBaseController) ParameterChecker4Page(page int) {
	if page <= 0 {
		logs.Error("[ParameterChecker4Page] page out of range, router: %s, traceID: %s, accountID: %d, ip: %s, reqs: %s, page: %d",
			c.CurrentRouter, c.TraceID, c.AccountID, c.IP, c.RequestData, page)
		c.TerminateWithCode(cerror.ParameterValueOutOfRange)
		return
	}
}

func (c *APIBaseController) ParameterChecker4Size(size int) {
	if size < 5 || size > 1000 {
		logs.Error("[ParameterChecker4Size] size out of range, router: %s, traceID: %s, accountID: %d, ip: %s, reqs: %s, size: %d",
			c.CurrentRouter, c.TraceID, c.AccountID, c.IP, c.RequestData, size)
		c.TerminateWithCode(cerror.ParameterValueOutOfRange)
		return
	}
}

func (c *APIBaseController) CommonResp(code cerror.ErrCode, data interface{}) {
	c.Data["json"] = c.BuildApiResponseEncrypt(code, data)
	c.ServeJSON()
}

func (c *APIBaseController) SuccessResp(data interface{}) {
	c.Data["json"] = c.BuildApiResponseEncrypt(cerror.CodeSuccess, data)
	c.ServeJSON()
}

func (c *APIBaseController) ServeJSONExt() {
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")

	content, err := libtools.JSONMarshal(c.Data["json"])
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = c.Ctx.Output.Body(content)
}

func (c *APIBaseController) TerminateWithCode(code cerror.ErrCode) {
	if c.XTrace {
		// 打印响应体主数据,以供联调排查问题
		logs.Notice("[trace] TerminateWithCode, router: %s, traceID: %s, ip: %s, accountID: %d, code: %d",
			c.CurrentRouter, c.TraceID, c.IP, c.AccountID, int(code))
	}

	c.Data["json"] = cerror.BuildApiResponse(code, "")
	c.ServeJSON()
	c.Abort("")
	return
}

func (c *APIBaseController) Finish() {
	// 量化接口性能
	duration := time.Since(c.beginTime)
	metrics.WebRequestDuration.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Observe(duration.Seconds())

	logs.Info("[Finish] client request completed, router: %s, traceID: %s, ip: %s, accountID: %d ",
		c.CurrentRouter, c.TraceID, c.IP, c.AccountID)
}

func (c *APIBaseController) BuildApiResponseEncrypt(code cerror.ErrCode, data interface{}) cerror.ApiResponse {
	r := cerror.ApiResponse{
		Code:      code,
		Message:   i18n.T(c.Lang, cerror.ErrorMessage(code)),
		SeverTime: libtools.GetUnixMillis(),
		Data:      data,
	}

	if code == cerror.CodeSuccess || code == cerror.AppForceUpgrade {
		//jsonByte, err := json.Marshal(data)
		jsonByte, err := libtools.JSONMarshal(data)
		if err != nil {
			r.Code = cerror.ServiceUnavailable
			r.Message = i18n.T(c.Lang, cerror.ErrorMessage(cerror.ServiceUnavailable))
			return r
		}

		encryptData, err := libtools.AesEncryptCBC(string(jsonByte), c.AesCBCKey, c.AesCBCIV)
		if err != nil {
			return r
		}

		r.Data = encryptData

		if c.XTrace {
			// 打印响应体主数据,以供联调排查问题
			logs.Notice("[trace] build output, router: %s, traceID: %s, ip: %s, accountID: %d, aesKey: %s, aesIV: %s, data: %s",
				c.CurrentRouter, c.TraceID, c.IP, c.AccountID, c.AesCBCKey, c.AesCBCIV, string(jsonByte))
		}
	} else {
		logs.Debug(">>> empty response, router: %s, traceID: %s, ip: %s, accountID: %d, code: %d",
			c.CurrentRouter, c.TraceID, c.IP, c.AccountID, int(code))
	}

	return r
}
