package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/cerror"
	"tinypro/common/models"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/account"
	"tinypro/common/pkg/system/config"
	"tinypro/common/pkg/tagslib"
	"tinypro/common/pkg/weixin"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebApiController struct {
	WebApiBaseController
}

func (c *WebApiController) Prepare() {
	// 调用上一级的 Prepare 方
	c.WebApiBaseController.Prepare()
}

func (c *WebApiController) Ping() {
	data := map[string]interface{}{
		"server_time": libtools.GetUnixMillis(),
		"version":     types.WebApiVersion,
		"head_hash":   libtools.GitRevParseHead(),
		"router":      c.Ctx.Request.RequestURI,
		"ip":          c.IP,
		"user_agent":  c.Ctx.Input.UserAgent(),
	}
	c.SuccessResponse(data)
}

func (c *WebApiController) IsLogin() {
	data := map[string]interface{}{
		"is_login":    1,
		"server_time": libtools.GetUnixMillis(),
	}

	c.SuccessResponse(data)
}

func (c *WebApiController) OauthLogin() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"open_oauth_plt": true,
		"open_user_id":   true,
		//"nickname":       true,
		//"open_avatar":    true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var loginReq types.ApiOauthLoginReqT
	err := json.Unmarshal(c.RequestBody, &loginReq)
	if err != nil || loginReq.OpenUserID == "" || loginReq.OpenOauthPlt < 1 {
		logs.Error("[OauthLogin] parse request get exception, ip: %s, accountID: %d, reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	// 简化入驻流程
	if loginReq.Nickname == "" {
		loginReq.Nickname = account.GenGuestNickname()
	}

	user, err := account.RegisterOrLogin(loginReq, c.IP, types.WebApiVersion)
	if err != nil || user.Id <= 0 {
		logs.Error("[OauthLogin] register or login get exception, ip: %s, accountID: %d, reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceUnavailable)
		return
	}

	accessToken, err := accesstoken.GenTokenWithCache(user.Id, types.PlatformWxMiniProgram, c.IP)
	if err != nil {
		logs.Error("[OauthLogin] gen token get exception, ip: %s, accountID: %d, reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceUnavailable)
		return
	}

	data := map[string]interface{}{
		"access_token": accessToken,
	}

	//if user.Mobile == "" {
	//	c.CommonResponse(cerror.NeedBindMobile, data)
	//} else {
	//	c.SuccessResponse(data)
	//}

	c.SuccessResponse(data)
}

func (c *WebApiController) GetConfig() {
	type resT struct {
		//TaskSteps string `json:"task_steps"`
		Risk struct {
			SpeedMax     int `json:"speed_max"`
			HeartRateMin int `json:"heart_rate_min"`
			HeartRateMax int `json:"heart_rate_max"`
		} `json:"risk"`
		CalorieConfigNum string `json:"calorie_config_num"`
		StepConfigNum    string `json:"step_config_num"`
	}

	var res resT

	speedMax := config.ValidItemString("risk_speed_max")
	heartRateMin := config.ValidItemString("risk_heart_rate_min")
	heartRateMax := config.ValidItemString("risk_heart_rate_max")

	res.Risk.SpeedMax, _ = libtools.Str2Int(speedMax)
	res.Risk.HeartRateMin, _ = libtools.Str2Int(heartRateMin)
	res.Risk.HeartRateMax, _ = libtools.Str2Int(heartRateMax)
	res.CalorieConfigNum = "0.65"
	res.StepConfigNum = "0.65"

	c.SuccessResponse(res)
}

func (c *WebApiController) Faq() {
	var res = make([]types.ApiFaqItem, 0)

	o := orm.NewOrm()
	m := models.Faq{}

	var list []models.Faq
	_, err := o.QueryTable(m.TableName()).
		Filter("status", types.StatusValid).
		OrderBy("weight", "id").Limit(10).All(&list)
	if err != nil {
		logs.Error("[Faq] db filter exception, ip: %s, accountID: %d, err: %v", c.IP, c.AccountID, err)
	}

	for _, data := range list {
		faq := types.ApiFaqItem{
			Subject: data.Subject,
			Content: data.Content,
			Tags:    tagslib.DataTagTupleGroup(data.Id),
		}

		res = append(res, faq)
	}

	c.SuccessResponse(res)
}

func (c *WebApiController) FaqMore() {
	var res = make([]types.ApiFaqItem, 0)

	// 必要参数检查
	requiredParameter := map[string]bool{
		"page": true,
		"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		Page int `json:"page"`
		Size int `json:"size"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[FaqMore] parse request get exception, ip: %s, accountID: %d reqs: %s, err: %v",
			c.IP, c.AccountID, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	c.ParameterChecker4Page(req.Page)
	c.ParameterChecker4Size(req.Size)

	o := orm.NewOrm()
	m := models.Faq{}

	offset := (req.Page - 1) * req.Size

	var list []models.Faq
	_, err = o.QueryTable(m.TableName()).
		Filter("status", types.StatusValid).
		OrderBy("weight", "id").
		Limit(req.Size).Offset(offset).
		All(&list)
	if err != nil {
		logs.Error("[FaqMore] db filter exception, ip: %s, accountID: %d, err: %v", c.IP, c.AccountID, err)
	}

	for _, data := range list {
		faq := types.ApiFaqItem{
			Subject: data.Subject,
			Content: data.Content,
			Tags:    tagslib.DataTagTupleGroup(data.Id),
		}

		res = append(res, faq)
	}

	c.SuccessResponse(res)
}

func (c *WebApiController) OpBroadcast() {
	var res []string

	var userObj models.AppUser
	var userList []models.AppUser
	_, _ = models.OrmList(&userObj, nil, 1, 10, false, &userList)
	for _, user := range userList {
		res = append(res, fmt.Sprintf(`欢迎新会员 %s 加入！%s`, user.Nickname, libtools.PickRandomEmoji()))
	}

	libtools.ShuffleStringList(res)

	c.SuccessResponse(res)
}

func (c *WebApiController) UpdateProfile() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"nickname":    true,
		"open_avatar": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		Nickname   string `json:"nickname"`
		OpenAvatar string `json:"open_avatar"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil || req.OpenAvatar == "" || req.Nickname == "" {
		logs.Error("[UpdateProfile] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	user := c.AppUser
	user.Nickname = req.Nickname
	user.OpenAvatar = req.OpenAvatar
	_, err = models.OrmUpdate(&user, []string{"Nickname", "OpenAvatar"})
	if err != nil {
		logs.Error("[UpdateProfile] update exception, ip: %s, accountId: %d, reqs: %s, err: %v", c.IP, c.AccountID, c.RequestData, err)
	} else {
		account.ResetAppUserNickname(c.AccountID, req.Nickname)
	}

	c.SuccessResponse(types.H{
		"op_msg": `操作成功`,
	})
}

func (c *WebApiController) WxDecrypt() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"raw":         true,
		"iv":          true,
		"session_key": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	type reqT struct {
		RawData    string `json:"raw"`
		IV         string `json:"iv"`
		SessionKey string `json:"session_key"`
	}

	var req reqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil || req.RawData == "" || req.IV == "" || req.SessionKey == "" {
		logs.Error("[WxDecrypt] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	text, err := weixin.Decrypt(req.RawData, req.SessionKey, req.IV)
	if err != nil {
		logs.Error("[WxDecrypt] decrypt get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidData)
		return
	}

	c.SuccessResponse(types.H{
		`text`: text,
	})
}
