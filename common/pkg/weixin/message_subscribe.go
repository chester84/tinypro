// 小程序消息订阅
// docs:
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
//
// 模板说明: https://shimo.im/docs/JYkgTGhv8XXxPgHW

package weixin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tinypro/common/models"
	"tinypro/common/types"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
)

type SubscribeMessageTplEnum string

const (
	SubEnrollCheck SubscribeMessageTplEnum = `enroll-check` //报名审核通知
	SubAttendClass SubscribeMessageTplEnum = `attend-class` //订阅课程开课提醒
)

type SubscribeMessageTplItem struct {
	DbSN SubscribeMessageTplEnum `json:"db_sn"`

	Name string `json:"name"`

	mnpTplId string // 小程序订阅模板编号,不导出
}

func SubscribeMessageTplConf() map[SubscribeMessageTplEnum]SubscribeMessageTplItem {
	return map[SubscribeMessageTplEnum]SubscribeMessageTplItem{
		SubEnrollCheck: {
			SubEnrollCheck,
			"报名审核通知",
			`h3IybQ9m6kqWtQOZN9WackqKQGfSLEDoQaFEg6iyokQ`,
		},
		SubAttendClass: {
			SubAttendClass,
			"订阅课程开课提醒",
			`UKp_zq3qN4H0tylWHQGcDFbah8UOyAbRzeC2EnCzReY`,
		},
	}
}

type SubscribeMessageSendDataItem struct {
	Value string `json:"value"`
}

type SubscribeMessageSendRequest struct {
	AccessToken string `json:"access_token"`

	ToUser     string `json:"touser"`
	TemplateId string `json:"template_id"`
	Page       string `json:"page"`
	Lang       string `json:"lang,omitempty"`

	Data map[string]SubscribeMessageSendDataItem `json:"data"`

	MiniProgramState string `json:"miniprogram_state,omitempty"`
}

func buildSendSubscribeMessageApi(appSN int) string {
	token, _ := getClientCredentialToken(appSN)
	return fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s`, token)
}

func SendSubscribeMessage(appSN int, userObj models.AppUser, courseObj models.Course, tplId SubscribeMessageTplEnum, page string, data map[string]SubscribeMessageSendDataItem) error {
	if !libtools.IsProductEnv() {
		//logs.Warning("[SendSubscribeMessage] 目前仅在生产环境可用")
		//return nil
	}

	api := buildSendSubscribeMessageApi(appSN)
	accessToken, _ := getClientCredentialToken(appSN)

	req := SubscribeMessageSendRequest{
		AccessToken: accessToken,
		ToUser:      userObj.WxOpenId,
		TemplateId:  SubscribeMessageTplConf()[tplId].mnpTplId,
		Page:        page,
		Data:        data,
	}
	reqJson, err := libtools.JsonEncode(req)
	if err != nil {
		logs.Error("[SendSubscribeMessage] can not json encode, reqs: %#v, err: %v", req, err)
		return err
	}

	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodPOST, api, nil, reqJson, libtools.DefaultHttpTimeout())

	log := models.WxInvokeLog{
		ApiType:   types.RemindAttendClass,
		UserID:    userObj.Id,
		CourseId:  courseObj.ID,
		Api:       api,
		Param:     reqJson,
		RespCode:  httpCode,
		Resp:      string(httpBody),
		CreatedAt: time.Now(),
	}
	models.OrmInsert(&log)

	if err != nil || httpCode != http.StatusOK {
		logs.Error("[SendSubscribeMessage] call api get exception, api: %s, reqs: %s, err: %v", api, reqJson, err)
		return err
	}

	logs.Debug("[SendSubscribeMessage] request: %s, response: %s", reqJson, string(httpBody))

	type resT struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
	}
	var res resT
	err = json.Unmarshal(httpBody, &res)
	if err != nil {
		logs.Error("[SentTplMessage] response data exception, reqs: %s, res: %s, err: %v", reqJson, string(httpBody), err)
		return err
	}

	if res.ErrCode != 0 {
		err = fmt.Errorf(`unexpected return code, response: %s`, string(httpBody))
		return err
	}

	return nil
}
