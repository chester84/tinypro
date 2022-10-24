// https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Overview.html

/*

https://shimo.im/docs/tXRvVXjJ6jrwyV9G/read

以最新文档为准
*/

package weixin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
)

type TemplateSendDataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type TplMiniProgram struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type TemplateSendRequest struct {
	ToUser     string `json:"touser"`
	TemplateId string `json:"template_id"`
	Url        string `json:"url,omitempty"`

	MiniProgram *TplMiniProgram `json:"miniprogram,omitempty"`

	Data map[string]TemplateSendDataItem `json:"data"`
}

const (
	// 公众号 -> 功能 -> 模板消息
	YcNewOrderNoticeTpl   = `50MDOcAFLN9H5hlrEIkORj5fxy_sCt8xFw2xV5OCM_w` // 新订单通知
	YcQuoteNoticeTpl      = `5TBu8p9U_kYBpXeSuMfPmtPKqQ1Nod7FQTjx0bH7-ig` // 报价成功通知
	YcQuoteFeedbackTpl    = `A415pd1zBmCwEF5WsgQzh_K_actCEJkj1KC0PNA8AX8` // 报价申请反馈通知
	YcProgressReminderTpl = `HhfQKxGKy8btAbueGqGe-6oFTvKvUQWrDZg3wO854SQ` // 受理进度提醒
	YcCaseReminderTpl     = `c6CM4Ri9OrBRbDyGgW5z51VIsHlIF_0S-Ar5n6SwQSE` // 理赔结案提醒
	YcSupplyNoticeTpl     = `rFKKdJx-TOFVXzSPcQ6TkGxVuUpYO04JayXQjFlHf3c` // 资料补充通知
	YcTransNoticeTpl      = `tc9-kFdA7GhfRmQnekT28cfCKTX-BwTNsSHNI7EUK3Y` // 交易提醒
	YcIncidentNoticeTpl   = `udNXAkluch6mS3bH0fEZCgtY4ORM0DbKY3fYhJdlZeE` // 事故进度提醒
	YcUpPhotoNoticeTpl    = `xqgLl65h_HKOwPMoIrt4Tw3wruM2SkSMbnoOfPD6R58` // 照片上传成功通知
	YcMutualSuccessTpl    = `RRPtJdd9qPdotIYQ60IaUkxPrpH_jtzJwF5nszu269w` // 互助成功通知

	CommentRemark = `点击查看详细报价信息。如有疑问，请联系在线客服。`
)

func CallbackMessageToken() string {
	if libtools.IsProductEnv() {
		return `654347eb6ef1c78622df94d1399a51e2`
	} else {
		return ``
	}
}

func CallbackMessageAESKey() string {
	if libtools.IsProductEnv() {
		return `9orALbYLUT44XBr61HdffADagGMPu4G1OSbIsXdEso3`
	} else {
		return ``
	}
}

func buildTplMessageApi(appSN int) string {
	token, _ := getClientCredentialToken(appSN)
	return fmt.Sprintf(`https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s`, token)
}

func SentTplMessage(appSN int, url string, openId string, tplId string, miniProgram *TplMiniProgram, data map[string]TemplateSendDataItem) (int64, error) {
	if !libtools.IsProductEnv() {
		logs.Warning("[SentTplMessage] 目前仅在生产环境可用")
		return 0, nil
	}

	api := buildTplMessageApi(appSN)

	req := TemplateSendRequest{
		ToUser:      openId,
		TemplateId:  tplId,
		Url:         url,
		MiniProgram: miniProgram,
		Data:        data,
	}
	reqJson, err := libtools.JsonEncode(req)
	if err != nil {
		logs.Error("[SentTplMessage] can not json encode, reqs: %#v, err: %v", req, err)
		return 0, err
	}

	httpBody, httpCode, err := libtools.SimpleHttpClient(libtools.HttpMethodPOST, api, nil, reqJson, libtools.DefaultHttpTimeout())
	if err != nil || httpCode != http.StatusOK {
		logs.Error("[SentTplMessage] call api get exception, api: %s, reqs: %s, err: %v", api, reqJson, err)
		return 0, err
	}

	logs.Debug("[SentTplMessage] request: %s, response: %s", reqJson, string(httpBody))

	type resT struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MsgId   int64  `json:"msgid"`
	}
	var res resT
	err = json.Unmarshal(httpBody, &res)
	if err != nil || res.ErrMsg != "ok" {
		logs.Error("[SentTplMessage] response data exception, reqs: %s, res: %s, err: %v", reqJson, string(httpBody), err)
		return 0, err
	}

	return res.MsgId, nil
}
