package controllers

import (
	"encoding/json"
	"tinypro/common/pkg/system/config"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus"

	"tinypro/common/cerror"
	"tinypro/common/models"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/metrics"
	"tinypro/common/pkg/security"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebApiBaseController struct {
	beego.Controller

	RequestBody []byte
	RequestData string

	CurrentRouter string // 当前路由

	XTrace      bool
	TraceID     string
	RequestTime int64

	// request json
	RequestJSON map[string]interface{}
	// 有效token对应的用户账户
	AccountID int64
	AppUser   models.AppUser
	IP        string

	beginTime time.Time
}

func (c *WebApiBaseController) BuildApiResponse(code cerror.ErrCode, data interface{}) cerror.ApiResponse {
	r := cerror.ApiResponse{
		Code:      code,
		Message:   cerror.ErrorMessage(code),
		SeverTime: libtools.GetUnixMillis(),
		Data:      data,
	}

	if c.XTrace || !libtools.IsProductEnv() {
		// 打印响应体主数据,以供联调排查问题
		jsonByte, _ := libtools.JSONMarshal(r)
		logs.Notice("[trace] build output, router: %s, traceID: %s, ip: %s, accountID: %d, data: %s",
			c.CurrentRouter, c.TraceID, c.IP, c.AccountID, string(jsonByte))
	}

	return r
}

func (c *WebApiBaseController) CommonResponse(code cerror.ErrCode, data interface{}) {
	c.Data["json"] = c.BuildApiResponse(code, data)
	c.ServeJSON()
}

func (c *WebApiBaseController) TerminateWithCode(code cerror.ErrCode) {
	c.Data["json"] = c.BuildApiResponse(code, cerror.EmptyData)
	c.ServeJSON()
	c.Abort("")
	return
}

func (c *WebApiBaseController) TerminateWithCodeAndData(code cerror.ErrCode, data interface{}) {
	c.Data["json"] = c.BuildApiResponse(code, data)
	c.ServeJSON()
	c.Abort("")
	return
}

func (c *WebApiBaseController) SuccessResponse(data interface{}) {
	c.Data["json"] = c.BuildApiResponse(cerror.CodeSuccess, data)
	c.ServeJSON()
}

func (c *WebApiBaseController) Prepare() {
	c.IP = c.Ctx.Input.IP()
	c.CurrentRouter = c.Ctx.Input.URL()

	// 量化接口并发量
	metrics.WebRequestTotal.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Inc()

	c.beginTime = time.Now()

	xTrace := c.Ctx.Request.Header.Get("X-Trace")
	if xTrace == types.XTrace {
		c.XTrace = true
	}

	c.RequestBody = c.Ctx.Input.RequestBody
	//data := c.GetString("data")
	c.RequestData = string(c.RequestBody)
	if len(c.RequestData) < 16 {
		logs.Warning("[base->Prepare] post data is empty.")
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	if !libtools.IsProductEnv() || c.XTrace {
		// 联调打印原始数据
		logs.Notice(">>> router: %s, ip: %s, RequestBody: %s", c.CurrentRouter, c.IP, c.RequestData)
	}

	var reqData types.WebApiBaseT
	err := json.Unmarshal(c.RequestBody, &reqData)
	if err != nil {
		logs.Warning("[base->Prepare] cat NOT json decode request body, router: %s, ip: %s, data: %s, err %#v", c.CurrentRouter, c.IP, string(c.RequestBody), err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	//白名单逻辑，等到功能测试完成
	//此逻辑需要注释掉
	//c.WhiteListCheck(reqData)

	requestTime, _ := libtools.Str2Int64(reqData.RequestTime)

	// 请求时间和服务器时间相差不能超过指定时间
	/*
		timeNow := libtools.GetUnixMillis()
		timeDiff := libtools.AbsInt64(timeNow - requestTime)
		if reqData.TraceID == "" ||
			timeDiff > 600000 {
			logs.Warning("[base->Prepare] capture suspected attacks, router: %s, traceID: %s, ip: %s, requestTime: %d, sysTime: %d, timeDiff: %d, httpBody: %s",
				c.CurrentRouter, c.TraceID, c.IP, reqData.RequestTime, timeNow, timeDiff, string(c.RequestBody))
			c.TerminateWithCode(cerror.ClientTimeASync)
			return
		}
	*/

	c.TraceID = reqData.TraceID
	c.RequestTime = requestTime
	logs.Info("[Prepare] client start request, router: %s, traceID: %s, ip: %s",
		c.CurrentRouter, c.TraceID, c.IP)

	// 抓包数据防重入
	if !security.PassStrongPreventRepeatedEntry(c.CurrentRouter, c.TraceID, c.IP) {
		logs.Warning("[WebApiBaseController] repeated entry, router: %s, traceID: %s, ip: %s, httpBody: %s",
			c.CurrentRouter, c.TraceID, c.IP, string(c.RequestBody))
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	var reqJSON map[string]interface{}
	err = json.Unmarshal(c.RequestBody, &reqJSON)
	if err != nil {
		logs.Warning("[base->Prepare] cat NOT json decode request data: %s, err %#v", c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	// json decode 通过
	c.RequestJSON = reqJSON

	// 必要参数检查,只检查存在,没有判值
	requiredParameter := map[string]bool{
		"trace_id":     true,
		"request_time": true,
		"access_token": true,
		"app_sn":       true,
	}
	//! 通用的必传参数在此外做统一校验
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		logs.Warning("[base->Prepare] lost required parameters, router: %s, traceID: %s, ip: %s, request: %s",
			c.CurrentRouter, c.TraceID, c.IP, c.RequestData)
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	uri := c.Ctx.Request.RequestURI
	// 以下路由不需要持有 token
	notNeedTokenRoute := map[string]bool{
		"/web-api/ping":              true,
		"/web-api/oauth-login":       true,
		"/web-api/mnp-check-version": true,
		"/web-api/get-config":        true,
		"/open-api/wx-oauth2-silent": true,

		`/web-api/edu/config`:          true,
		`/web-api/edu/search`:          true,
		`/web-ap/edu/login/send-sms`:   true,
		`/web-ap/edu/login/verify-sms`: true,
		`/web-ap/edu/login/sid-pwd`:    true,
		"/web-api/edu/landing-page":    true,
		"/web/api/edu/content-list":    true,
		`/web/api/edu/content-detail`:  true,
		"/web-api/edu/hot-zone-map":    true,
		`/web-api/edu/apply`:           true,

		`/web-api/edu/topic-content-list`: true,

		// 微信解密接口
		`/web-api/wx/decrypt`: true,

		// 四维相关
		`/web-api/config`:            true,
		`/web-api/wx/login-register`: true,
		`/web-api/front-page`:        true,
		`/web-api/public-courses`:    true,
	}
	if !notNeedTokenRoute[uri] {
		// 检查 token 有效性
		var (
			ok        bool
			ok2       bool
			accountId int64
			token     = reqJSON["access_token"].(string)
		)

		ok, accountId = accesstoken.IsValidAccessToken(types.PlatformWxMiniProgram, token)
		if !ok {
			ok2, accountId = accesstoken.IsValidAccessToken(types.PlatformH5, token)
			if !ok2 {
				logs.Notice("[base->Prepare] access_token is invalid, json:", c.RequestData)
				c.TerminateWithCode(cerror.InvalidAccessToken)
				return
			}
		}

		c.AccountID = accountId
		err = models.OrmOneByPkId(accountId, &c.AppUser)
		if err != nil {
			logs.Notice("[base->Prepare] user is empty, accountId is %d", accountId)
			c.TerminateWithCode(cerror.InvalidAccessToken)
			return
		}
	}

	c.RequestJSON["ip"] = c.Ctx.Input.IP()
}

func (c *WebApiBaseController) WhiteListCheck(reqData types.WebApiBaseT) {

	logs.Debug("WhiteListCheck %#v", reqData)
	//白名单imei, 此功能维护后后，要去掉
	if reqData.Imei != "" {
		c.InWhiteList(reqData.Imei)
	} else {
		// 检查 token 有效性
		tokenObj, err := accesstoken.GetUserIdByToken(reqData.AccessToken)
		if err != nil {
			logs.Error("[whiteList] cant get AccountToken reqData %#v", reqData)
			c.TerminateWithCode(cerror.ServiceUnavailable)
			return
		}

		var userObj models.AppUser
		err = models.OrmOneByPkId(tokenObj.AccountId, &userObj)
		if err != nil {
			logs.Error("[whiteList] User cant be found reqs %v", reqData)
			c.TerminateWithCode(cerror.ServiceUnavailable)
			return
		}

		//todo 目前看新手表用户，被此处逻辑限制了，先去掉

		//var watchObj models.WatchDevice
		//err = models.OrmOneByPkId(userObj.WatchDeviceId, &watchObj)
		//if err != nil {
		//	logs.Error("[whiteList] watchObj cant be found reqs %#v, user %#v", reqData, userObj)
		//	c.TerminateWithCode(cerror.ServiceUnavailable)
		//	return
		//}
		//
		//c.InWhiteList(watchObj.IMEI)
	}
}

func (c *WebApiBaseController) InWhiteList(imei string) {
	if !inWhiteList(imei) {
		// 停机维护公告
		logs.Debug("%s not in whitelist", imei)
		var stop4Upgrade bool
		stop4Upgrade = true
		if stop4Upgrade {
			c.CommonResponse(cerror.Stop4Upgrade, cerror.EmptyData)
			return
		}
	} else {
		logs.Debug("%s in whitelist", imei)
	}
}

func inWhiteList(imei string) bool {
	whiteImei := config.ValidItemString("white-list-imei")
	var whiteImeiList []string
	logs.Debug("whiteImei %#v", whiteImei)
	err := json.Unmarshal([]byte(whiteImei), &whiteImeiList)
	if err != nil {
		logs.Error("[whiteImeiList] cat NOT json decode, err %#v", err)
	}

	if libtools.InSlice(imei, whiteImeiList) {
		return true
	}

	return false
}

func (c *WebApiBaseController) Finish() {
	// 量化接口性能
	duration := time.Since(c.beginTime)
	metrics.WebRequestDuration.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Observe(duration.Seconds())
}

func (c *WebApiBaseController) ParameterChecker4Size(size int) {
	if size < 5 || size > 100 {
		logs.Error("[ParameterChecker4Size] size out of range, router: %s, traceID: %s, accountID: %d, ip: %s, reqs: %s, size: %d",
			c.CurrentRouter, c.TraceID, c.AccountID, c.IP, c.RequestData, size)
		c.TerminateWithCode(cerror.ParameterValueOutOfRange)
		return
	}
}

func (c *WebApiBaseController) ParameterChecker4Page(page int) {
	if page <= 0 {
		logs.Error("[ParameterChecker4Page] page out of range, router: %s, traceID: %s, accountID: %d, ip: %s, reqs: %s, size: %d",
			c.CurrentRouter, c.TraceID, c.AccountID, c.IP, c.RequestData, page)
		c.TerminateWithCode(cerror.ParameterValueOutOfRange)
		return
	}
}
