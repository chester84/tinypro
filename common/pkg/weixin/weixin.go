// docs: https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
package weixin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"

	"tinypro/common/lib/redis/cache"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type Oauth2SilentResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

type SnsApiUserInfo struct {
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	Unionid    string `json:"unionid"`

	Sex types.GenderEnum `json:"sex"` // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知

	//Privilege  []string `json:"privilege"`
}

func (r *SnsApiUserInfo) FixWxSex() types.GenderEnum {
	var gender = types.GenderUnknown
	switch r.Sex {
	case 1:
		gender = types.GenderMale

	case 2:
		gender = types.GenderFemale
	}

	return gender
}

type SnsJsCode2SessionResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int64  `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

type ClientCredentialToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
}

type WxJsApiConfig struct {
	Debug     bool     `json:"debug"`
	AppId     string   `json:"appId"`
	Timestamp int64    `json:"timestamp"`
	NonceStr  string   `json:"nonceStr"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

type TicketResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

/*
签约回调通知url
*/
func EntrustNotifyUrl() string {
	url := fmt.Sprintf("%s%s", libtools.InternalApiDomain(), "/wx/entrust")
	return url
}

func Oauth2Silent(code string, appSN int) (response Oauth2SilentResponse, err error) {
	var appid = AppID(appSN)
	var secret = Secret(appSN)

	var api = fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		appid, secret, code)

	reqHeaders := map[string]string{}
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, reqHeaders, "", libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[Oauth2Silent] call weixin oauth2 fail, api: %s, err: %v", api, err)
		return
	}

	logs.Warning("[Oauth2Silent] httpCode: %d, response: %s, err: %v", httpCode, string(httpBody), err)

	err = json.Unmarshal(httpBody, &response)
	if err != nil {
		logs.Error("[Oauth2Silent] json decode exception, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	return
}

func SnsUserInfo(accessToken, openid string) (info SnsApiUserInfo, err error) {
	var api = fmt.Sprintf(`https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN`, accessToken, openid)
	reqHeaders := map[string]string{}

	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, reqHeaders, "", libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[SnsUserInfo] call weixin sns/userinfo fail, api: %s, err: %v", api, err)
		return
	}

	logs.Warning("[SnsUserInfo] httpCode: %d, response: %s, err: %v", httpCode, string(httpBody), err)

	err = json.Unmarshal(httpBody, &info)
	if err != nil {
		logs.Error("[SnsUserInfo] json decode exception, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	return
}

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func SnsJsCode2Session(code string, appSN int) (response SnsJsCode2SessionResponse, err error) {
	var api = fmt.Sprintf(`https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code`,
		AppID(appSN), Secret(appSN), code)

	reqHeaders := map[string]string{}
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, reqHeaders, "", libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[SnsJsCode2Session] call weixin oauth2 fail, api: %s, err: %v", api, err)
		return
	}

	logs.Warning("[SnsJsCode2Session] httpCode: %d, response: %s, err: %v", httpCode, string(httpBody), err)

	err = json.Unmarshal(httpBody, &response)
	if err != nil {
		logs.Error("[SnsJsCode2Session] json decode exception, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	if response.Openid == "" && response.Unionid == "" {
		err = fmt.Errorf(`call wx api get exception, response: %s`, string(httpBody))
	}

	return
}

// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func getClientCredentialToken(appSN int) (token string, err error) {
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cKey := fmt.Sprintf(`%s:%d`, rdsKeyClientCredentialToken, appSN)
	cValue, err := redis.String(cacheClient.Do("GET", cKey))
	if err != nil && err != redis.ErrNil {
		logs.Error("[getClientCredentialToken] redis> GET %s, err: %v", cKey, err)
	}

	if cValue != "" {
		token = cValue
		return
	}

	var expires = 7000

	var response ClientCredentialToken
	api := fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s`,
		AppID(appSN), Secret(appSN))
	reqHeaders := map[string]string{}
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, reqHeaders, "", libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[getClientCredentialToken] call weixin token fail, api: %s, err: %v", api, err)
		return
	}

	logs.Warning("[getClientCredentialToken] httpCode: %d, response: %s, err: %v", httpCode, string(httpBody), err)

	err = json.Unmarshal(httpBody, &response)
	if err != nil {
		logs.Error("[getClientCredentialToken] json decode exception, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	if response.AccessToken == "" {
		err = fmt.Errorf(`call weixin api get empty token`)
		logs.Error("[getClientCredentialToken] call api fail, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	token = response.AccessToken
	_, err = cacheClient.Do("SETEX", cKey, expires, token)
	if err != nil {
		logs.Error("[getClientCredentialToken] redis> SETEX %s %d %s ,err: %v",
			cKey, expires, token)
	}

	return
}

func getJsApiTicket(appSN int) (ticket string) {
	cacheClient := cache.RedisCacheClient.Get()
	defer cacheClient.Close()

	cKey := fmt.Sprintf(`%s:%d`, rdsKeyJsApiTicket, appSN)
	cValue, err := redis.String(cacheClient.Do("GET", cKey))
	if err != nil && err != redis.ErrNil {
		logs.Error("[getJsApiTicket] redis> GET %s, err: %v", cKey, err)
	}
	if cValue != "" {
		ticket = cValue
		return
	}

	var expires = 7000

	accessToken, _ := getClientCredentialToken(appSN)
	var response TicketResponse
	api := fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi`, accessToken)
	reqHeaders := map[string]string{}
	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, api, reqHeaders, "", libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[getJsApiTicket] call weixin get ticket fail, api: %s, err: %v", api, err)
		return
	}

	logs.Warning("[getJsApiTicket] httpCode: %d, response: %s, err: %v", httpCode, string(httpBody), err)

	err = json.Unmarshal(httpBody, &response)
	if err != nil {
		logs.Error("[getJsApiTicket] json decode exception, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	if response.Ticket == "" {
		err = fmt.Errorf(`call weixin api get empty ticket`)
		logs.Error("[getJsApiTicket] call api fail, api: %s, response: %s, err: %v",
			api, string(httpBody), err)
		return
	}

	ticket = response.Ticket
	_, err = cacheClient.Do("SETEX", cKey, expires, ticket)
	if err != nil {
		logs.Error("[getJsApiTicket] redis> SETEX %s %d %s, err: %v", cKey, expires, ticket, err)
	}

	return
}

// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
func GetWxJsApiConfig(currentUrl string, appSN int) (cfg WxJsApiConfig) {
	if !libtools.IsProductEnv() {
		cfg.Debug = true
	}
	cfg.AppId = AppID(appSN)
	cfg.Timestamp = time.Now().Unix()
	cfg.NonceStr = libtools.GenerateRandomStr(16)
	cfg.JsApiList = []string{
		"updateAppMessageShareData",
		"updateTimelineShareData",
		"startRecord",
		"stopRecord",
		"onVoiceRecordEnd",
		"playVoice",
		"pauseVoice",
		"stopVoice",
		"onVoicePlayEnd",
		"uploadVoice",
		"downloadVoice",
		"chooseImage",
		"previewImage",
		"uploadImage",
		"downloadImage",
		"translateVoice",
		"getNetworkType",
		"openLocation",
		"getLocation",
		"hideOptionMenu",
		"showOptionMenu",
		"hideMenuItems",
		"showMenuItems",
		"hideAllNonBaseMenuItem",
		"showAllNonBaseMenuItem",
		"closeWindow",
		"scanQRCode",
		"chooseWXPay",
		"openProductSpecificView",
		"addCard",
		"chooseCard",
		"openCard",
	}

	ticket := getJsApiTicket(appSN)

	origin := fmt.Sprintf(`jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s`, ticket, cfg.NonceStr, cfg.Timestamp, currentUrl)
	logs.Debug("origin: %s", origin)
	signature := libtools.Sha1(origin)
	logs.Debug("signature: %s", signature)

	cfg.Signature = signature

	return
}
