package controllers

import (
	"encoding/json"
	"tinypro/common/models"
	"tinypro/common/pkg/advertisingpopup"
	"tinypro/common/pkg/course"
	"tinypro/common/pkg/enroll"
	"tinypro/common/pkg/msg"
	"tinypro/common/pkg/satissurvey"
	"tinypro/common/pogo/reqs"
	"tinypro/common/service/enrollbiz"
	"tinypro/common/service/loginbiz"
	"tinypro/common/service/subscribe_template_biz"
	"tinypro/common/types"

	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/cerror"
	"github.com/chester84/libtools"
)

type WebApiSiWeiController struct {
	WebApiBaseController
}

func (c *WebApiSiWeiController) Prepare() {
	// 调用上一级的 Prepare 方
	c.WebApiBaseController.Prepare()
}

func (c *WebApiSiWeiController) WxOauthLoginOrRegister() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"code":   true,
		"app_sn": true,
		//"reginfo": false,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.WxLoginReqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[WxOauth2SilentLogin] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	user, wxSession, token, err := loginbiz.WxOauthLoginOrRegister(req, c.IP)
	if err != nil {
		logs.Error("[WxOauth2SilentLogin] loginbiz.WxOauth2SlientLogin get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceUnavailable)
		return
	}

	data := map[string]interface{}{}

	data["customer_service_qr"] = types.CustomerServiceQr
	data["advertise"] = types.SuccessEmptyResp()
	advertise, err := advertisingpopup.GetUserAccessAdvertisingPopup(user, 1)
	if err != nil {
		logs.Error("WxOauth2SilentLogin GetUserAccessAdvertisingPopup get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	data["advertise"] = advertise

	if user.Id <= 0 {
		data["wx_session"] = wxSession
		data["login_info"] = ""
		c.TerminateWithCodeAndData(cerror.UserUnRegister, data)
		return
	}

	var loginUser types.ApiLoginUserInfoResponse
	loginUser.Nickname = user.Nickname
	loginUser.OpenAvatar = user.OpenAvatar
	loginUser.AccessToken = token

	data["wx_session"] = wxSession
	data["login_info"] = loginUser

	c.SuccessResponse(data)
}

func (c *WebApiSiWeiController) Config() {
	ret := map[string]interface{}{}
	ret["satis_survey_questions"] = satissurvey.SatisQuestion()
	ret["survey_score_config"] = satissurvey.SurveyScoreConfig()
	ret["customer_service_qr"] = types.CustomerServiceQr
	c.SuccessResponse(ret)
}

func (c *WebApiSiWeiController) FrontPage() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageSelectedInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[FrontPage] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}
	list, err := course.FrontPage(req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) PublicCourses() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[PublicCourses] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	list, err := course.PublicCourses(req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) GetLastEnroll() {
	requiredParameter := map[string]bool{
		"course_sn": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.EnrollReqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[Enroll] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	courseID, _ := libtools.Str2Int64(req.CourseSN)
	obj, err := enrollbiz.GetLastEnroll(c.AppUser, courseID)
	if err != nil {
		logs.Error("[Enroll] EnrollCourse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(obj)
}

func (c *WebApiSiWeiController) IsPassEnroll() {
	requiredParameter := map[string]bool{
		"course_sn": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.EnrollReqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("IsPassEnroll parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	courseID, _ := libtools.Str2Int64(req.CourseSN)

	m := models.Course{}
	err = models.OrmOneByPkId(courseID, &m)
	if err != nil {
		logs.Error("IsPassEnroll OrmOneByPkId Course get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	isPass, err := enrollbiz.IsPassEnroll(c.AppUser.Id, courseID)
	if err != nil {
		logs.Error("IsPassEnroll EnrollCourse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	data := map[string]interface{}{}
	data["result"] = isPass
	data["work_wx_qr"] = m.WorkWxQR

	c.SuccessResponse(data)
}

func (c *WebApiSiWeiController) Enroll() {
	requiredParameter := map[string]bool{
		"course_sn": true,
		"real_name": true,
		"mobile":    true,
		"company":   true,
		"position":  true,
		"residence": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.EnrollReqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[Enroll] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	err = enrollbiz.EnrollCourse(req, c.AccountID)
	if err != nil {
		logs.Error("[Enroll] EnrollCourse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(struct{}{})
}

func (c *WebApiSiWeiController) MyEnrolls() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[MyEnrolls] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	list, err := enroll.MyEnrolls(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) MyCourses() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[MyCourses] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	list, err := enroll.MyCourses(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) GetEnrollDetail() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"course_sn": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.CourseDetailReq
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[GetEnrollDetail] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	courseId, _ := libtools.Str2Int64(req.CourseSN)
	ret, err := enroll.GetEnrollDetail(c.AppUser.Id, courseId)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(ret)
}

func (c *WebApiSiWeiController) SignInCourse() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"course_sn": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.SignInReqT
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[SignInCourse] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	_, err = enroll.SignInCourse(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(types.SuccessEmptyResp())
}

func (c *WebApiSiWeiController) MyHistoryCourses() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[MyHistoryCourses] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	list, err := enroll.MyHistoryCourses(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) MyMsgList() {
	// 必要参数检查
	requiredParameter := map[string]bool{
		"type": true,
		//"size": true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.PageInfo
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[MyMsgList] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	list, err := msg.MyMsgList(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}
	c.SuccessResponse(list)
}

func (c *WebApiSiWeiController) SatisSurvey() {
	requiredParameter := map[string]bool{
		"course_sn": true,
		"q1":        true,
		"q2":        true,
		"q3":        true,
		"q4":        true,
		"q5":        true,
		"q6":        true,
		"q7":        true,
		"q8":        true,
		"q9":        true,
		"q10":       true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.QScore
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[SatisSurvey] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	err = satissurvey.UserGiveScore(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(types.SuccessEmptyResp())
}

func (c *WebApiSiWeiController) SubscribeTemplate() {
	requiredParameter := map[string]bool{
		"course_sn": true,
		"ids":       true,
	}
	if !libtools.CheckRequiredParameter(c.RequestJSON, requiredParameter) {
		c.TerminateWithCode(cerror.LostRequiredParameters)
		return
	}

	var req reqs.SubscribeTmpl
	err := json.Unmarshal(c.RequestBody, &req)
	if err != nil {
		logs.Error("[SubscribeTemplate] parse request get exception, ip: %s, reqs: %s, err: %v",
			c.IP, c.RequestData, err)
		c.TerminateWithCode(cerror.InvalidRequestData)
		return
	}

	err = subscribe_template_biz.SubscribeTemplate(c.AppUser, req)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(types.SuccessEmptyResp())
}

func (c *WebApiSiWeiController) AdvertisePopup() {
	url, err := advertisingpopup.GetUserAccessAdvertisingPopup(c.AppUser, 1)
	if err != nil {
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse(url)
}
