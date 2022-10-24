package controllers

import (
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus"

	"tinypro/common/cerror"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/account"
	"tinypro/common/pkg/metrics"
	"tinypro/common/pkg/weixin"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type OpenApiController struct {
	beego.Controller

	RequestBody []byte
	RequestData string
	// request json
	RequestJSON map[string]interface{}

	CurrentRouter string // 当前路由
	IP            string

	beginTime time.Time
}

func (c *OpenApiController) Prepare() {
	c.IP = c.Ctx.Input.IP()
	c.CurrentRouter = c.Ctx.Input.URL()

	// 量化接口并发量
	metrics.WebRequestTotal.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Inc()

	c.beginTime = time.Now()

	c.RequestBody = c.Ctx.Input.RequestBody
	c.RequestData = string(c.RequestBody)
	//if len(c.RequestData) < 16 {
	//	logs.Warning("post data is empty.")
	//	c.TerminateWithCode(cerror.LostRequiredParameters)
	//	return
	//}

	var reqJSON map[string]interface{}
	err := json.Unmarshal(c.RequestBody, &reqJSON)
	if err != nil {
		logs.Info("cat NOT json decode request data:", c.RequestData)
	} else {
		// json decode 通过
		c.RequestJSON = reqJSON
	}

	if !libtools.IsProductEnv() {
		// 联调打印原始数据
		logs.Notice(">>> router: %s, ip: %s, RequestBody: %s", c.CurrentRouter, c.IP, c.RequestData)
	}
}

func (c *OpenApiController) BuildApiResponse(code cerror.ErrCode, data interface{}) cerror.ApiResponse {
	r := cerror.ApiResponse{
		Code:      code,
		Message:   cerror.ErrorMessage(code),
		SeverTime: libtools.GetUnixMillis(),
		Data:      data,
	}

	if !libtools.IsProductEnv() {
		// 打印响应体主数据,以供联调排查问题
		jsonByte, _ := libtools.JSONMarshal(r)
		logs.Notice("[trace] build output, router: %s, ip: %s,  data: %s",
			c.CurrentRouter, c.IP, string(jsonByte))
	}

	return r
}

func (c *OpenApiController) CommonResponse(code cerror.ErrCode, data interface{}) {
	c.Data["json"] = c.BuildApiResponse(code, data)
	c.ServeJSON()
}

func (c *OpenApiController) TerminateWithCode(code cerror.ErrCode) {
	c.Data["json"] = c.BuildApiResponse(code, cerror.EmptyData)
	c.ServeJSON()
	c.Abort("")
	return
}

func (c *OpenApiController) SuccessResponse(data interface{}) {
	c.Data["json"] = c.BuildApiResponse(cerror.CodeSuccess, data)
	c.ServeJSON()
}

func (c *OpenApiController) Finish() {
	// 量化接口性能
	duration := time.Since(c.beginTime)
	metrics.WebRequestDuration.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Observe(duration.Seconds())
}

func (c *OpenApiController) WxOauth2Silent() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"code":   true,
		"app_sn": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		Code  string `json:"code"`
		AppSN int    `json:"app_sn"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[WxOauth2Silent] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	var (
		authData   types.ApiOauthLoginReqT
		sessionKey string
	)

	if req.AppSN == weixin.AppSNWxGzh {
		wxOauth2Res, err := weixin.Oauth2Silent(req.Code, req.AppSN)
		if err != nil || wxOauth2Res.Openid == "" {
			logs.Error("[WxOauth2Silent] ip: %s, err: %v", c.IP, err)
			c.TerminateWithCode(cerror.WeixinOauth2Fail)
			return
		}

		wxUserInfo, err := weixin.SnsUserInfo(wxOauth2Res.AccessToken, wxOauth2Res.Openid)
		if err != nil {
			logs.Error("[SnsUserInfo] ip: %s, err: %v", c.IP, err)
			c.TerminateWithCode(cerror.WeixinOauth2Fail)
			return
		}

		authData = types.ApiOauthLoginReqT{
			AppSN:        req.AppSN,
			OpenOauthPlt: types.OpenOauthWeChat,
			Nickname:     wxUserInfo.Nickname,
			OpenUserID:   wxUserInfo.Unionid,
			OpenAvatar:   wxUserInfo.HeadImgUrl,
			Gender:       wxUserInfo.FixWxSex(),
			WxOpenId:     wxUserInfo.Openid,
			Country:      wxUserInfo.Country,
			Province:     wxUserInfo.Province,
			City:         wxUserInfo.City,
		}
	} else {
		wxJsSession, err := weixin.SnsJsCode2Session(req.Code, req.AppSN)
		if err != nil {
			logs.Error("[SnsJsCode2Session] ip: %s, reqs: %s, err: %v", c.IP, c.RequestData, err)
			c.TerminateWithCode(cerror.WeixinOauth2Fail)
			return
		}

		// 兼容 unionid 为空的情况
		if wxJsSession.Unionid == "" {
			wxJsSession.Unionid = wxJsSession.Openid
		}

		authData = types.ApiOauthLoginReqT{
			AppSN:        req.AppSN,
			OpenOauthPlt: types.OpenOauthWeChat,
			Nickname:     account.GenGuestNickname(),
			OpenUserID:   wxJsSession.Unionid,
			WxOpenId:     wxJsSession.Openid,
			Gender:       types.GenderUnknown,
		}

		sessionKey = wxJsSession.SessionKey
	}

	user, err := account.RegisterOrLogin(authData, c.IP, types.WebApiVersion)
	if err != nil || user.Id <= 0 {
		logs.Error("[WxOauth2Silent] register or login get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceUnavailable)
		return
	}

	accessToken, err := accesstoken.GenTokenWithCache(user.Id, types.PlatformH5, c.IP)
	if err != nil {
		logs.Error("[WxOauth2Silent] gen token get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceUnavailable)
		return
	}

	var hasBaseInfo int
	if user.Nickname != "" && user.OpenAvatar != "" {
		hasBaseInfo = 1
	}

	data := map[string]interface{}{
		"access_token":  accessToken,
		"session_key":   sessionKey,
		"has_base_info": hasBaseInfo,
	}

	c.SuccessResponse(data)
}

func (c *OpenApiController) WxJsConfig() {
	type reqT struct {
		ShareUrl string `json:"share_url"`
		AppSN    int    `json:"app_sn"`
		SceneNum int    `json:"scene_num"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[WxJsConfig] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	if req.ShareUrl == "" {
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	wxJsConfig := weixin.GetWxJsApiConfig(req.ShareUrl, req.AppSN)

	type resT struct {
		WxJsConfig weixin.WxJsApiConfig       `json:"wx_js_config"`
		ShareData  weixin.AppMessageShareData `json:"share_data"`
	}

	var res resT

	res.WxJsConfig = wxJsConfig
	res.ShareData = weixin.BuildShareData(req.SceneNum)

	c.SuccessResponse(res)
}
