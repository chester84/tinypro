package controllers

import (
	"image"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tencentyun/cos-go-sdk-v5"

	"tinypro/common/cerror"
	"tinypro/common/models"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/metrics"
	"tinypro/common/pkg/rbacv2"
	"tinypro/common/pkg/tc"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebAdmBaseController struct {
	beego.Controller

	IP            string
	XTrace        bool
	CurrentRouter string // 当前路由
	IsHackDev     bool

	AccountId int64
	Account   models.Admin

	beginTime time.Time
}

var noCheckPolicyRouter = map[string]bool{
	"/web-adm/ping": true,

	"/web-adm/user/login":  true,
	"/web-adm/user/logout": true,

	"/web-adm/system-config":        true,
	"/web-adm/user/base-info":       true,
	"/web-adm/user/change-password": true,
	`/web-adm/edu/apply-classify`:   true,

	// 特殊路由,让小程序可以游客上传资源
	`/web-api/edu/upload-resource`: true,

	"/web-adm/upload-resource4tinymce": true,
}

func (c *WebAdmBaseController) Prepare() {
	c.IP = c.Ctx.Input.IP()
	c.CurrentRouter = c.Ctx.Input.URL()

	// 量化接口并发量
	metrics.WebRequestTotal.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Inc()

	c.beginTime = time.Now()

	xTrace := c.GetString(`x_trace`)
	if xTrace != "" {
		c.XTrace = true
	}

	accessToken := c.Ctx.Request.Header.Get("X-Access-Token")
	uri := c.Ctx.Request.RequestURI
	// 以下路由不需要持有 token
	notNeedTokenRoute := map[string]bool{
		`/web-adm/ping`:          true,
		`/web-adm/user/login`:    true,
		`/web-adm/login-captcha`: true,

		`/web-adm/upload-resource4tinymce`: true,
	}
	if !notNeedTokenRoute[uri] {
		// 检查 token 有效性
		ok, accountId := accesstoken.IsValidAccessToken(types.PlatformAdm, accessToken)
		if !ok {
			c.TerminateWithCode(cerror.AccessDenied)
			return
		}

		c.AccountId = accountId

		var adm models.Admin
		err := models.OrmOneByPkId(accountId, &adm)
		if err != nil {
			logs.Error("[Prepare] get wrong entry. ip: %s, accountId: %d, err: %v", c.IP, accountId, err)
			c.TerminateWithCode(cerror.AccessDenied)
			return
		}

		c.Account = adm
	}

	if !noCheckPolicyRouter[c.CurrentRouter] {
		if !rbacv2.HasPolicyByIdRouter(libtools.Int642Str(c.AccountId), c.CurrentRouter) {
			logs.Warn("[HasPolicyByIdRouter] false uid:%d url:%v", c.AccountId, c.CurrentRouter)
			c.Data["json"] = cerror.BuildAjaxResponse(cerror.AccessDenied, cerror.EmptyData)
			c.ServeJSON()
			return
		}
	}
}

func (c *WebAdmBaseController) ParseDateRangeUnlimited(timeField string, condBox map[string]interface{}) (timeStart, timeEnd int64) {
	condBox["_time_field_"] = timeField

	timeStart = libtools.Date2UnixMsec(c.GetString(`queryDateBegin`), `Y-m-d`)
	timeEnd = libtools.Date2UnixMsec(c.GetString(`queryDateEnd`), `Y-m-d`) + 3600*24*1000
	if timeStart > 0 && timeEnd > timeStart {
		condBox["ctime_start_time"] = timeStart
		condBox["ctime_end_time"] = timeEnd
	}

	if c.IsHackDev {
		condBox["_is_hack_dev_"] = c.IsHackDev
	}

	return
}

func (c *WebAdmBaseController) ParseDateRangeCommon(timeField string, condBox map[string]interface{}, enforce bool) (timeStart, timeEnd int64) {
	condBox["_time_field_"] = timeField

	var rangeLimit int64 = types.DayMillisecond * 8

	timeStart = libtools.Date2UnixMsec(c.GetString(`queryDateBegin`), `Y-m-d`)
	timeEnd = libtools.Date2UnixMsec(c.GetString(`queryDateEnd`), `Y-m-d`) + 3600*24*1000
	if timeStart > 0 && timeEnd > timeStart {
		if timeEnd-timeStart > rangeLimit && !c.IsHackDev {
			timeEnd = timeStart + rangeLimit
		}
		condBox["ctime_start_time"] = timeStart
		condBox["ctime_end_time"] = timeEnd
	}

	if enforce && timeStart <= 0 {
		logs.Debug("[ParseDateRangeCommon] get empty input")
		condBox["ctime_start_time"] = libtools.NaturalDay(-7)
		condBox["ctime_end_time"] = libtools.NaturalDay(0) + types.DayMillisecond
	}

	if c.IsHackDev {
		condBox["_is_hack_dev_"] = c.IsHackDev
	}

	return
}

func (c *WebAdmBaseController) BuildApiResponse(code cerror.ErrCode, data interface{}) cerror.ApiResponse {
	r := cerror.ApiResponse{
		Code:      code,
		Message:   cerror.ErrorMessage(code),
		SeverTime: libtools.GetUnixMillis(),
		Data:      data,
	}

	if c.XTrace || !libtools.IsProductEnv() {
		// 打印响应体主数据,以供联调排查问题
		jsonByte, _ := libtools.JSONMarshal(r)
		logs.Notice("[trace] build output, router: %s, ip: %s, accountID: %d, data: %s",
			c.CurrentRouter, c.IP, c.AccountId, string(jsonByte))
	}

	return r
}

func (c *WebAdmBaseController) CommonResponse(code cerror.ErrCode, data interface{}) {
	c.Data["json"] = c.BuildApiResponse(code, data)
	c.ServeJSON()
}

func (c *WebAdmBaseController) TerminateWithCode(code cerror.ErrCode) {
	c.Data["json"] = c.BuildApiResponse(code, cerror.EmptyData)
	c.ServeJSON()
	c.Abort("")
	return
}

func (c *WebAdmBaseController) SuccessResponse(data interface{}) {
	c.Data["json"] = c.BuildApiResponse(cerror.CodeSuccess, data)
	c.ServeJSON()
}

func (c *WebAdmBaseController) Finish() {
	// 量化接口性能
	duration := time.Since(c.beginTime)
	metrics.WebRequestDuration.With(prometheus.Labels{
		"hostname": metrics.Hostname,
		"url":      c.Ctx.Input.URL(),
	}).Observe(duration.Seconds())
}

func (c *WebAdmBaseController) uploadCore(f multipart.File, h *multipart.FileHeader, upFilename string, access int) (s3Key, fileMd5 string, code cerror.ErrCode, mediaInfo types.MediaSimpleInfo, err error) {
	fileBytes, _ := ioutil.ReadAll(f)
	fileMd5 = libtools.Md5Bytes(fileBytes)

	extension, mime, _ := libtools.DetectFileByteType(fileBytes)
	_, s3Key = libtools.BuildHashName(fileMd5, extension)
	mediaInfo.MediaType = mime
	if strings.Contains(mime, "image") {
		_, _ = f.Seek(0, 0)
		imgObj, _, errD := image.DecodeConfig(f)
		if errD != nil {
			logs.Warning("[uploadCore] img decode exception, err: %v", errD)
		} else {
			mediaInfo.Width = imgObj.Width
			mediaInfo.Height = imgObj.Height
		}
	}

	_, _ = f.Seek(0, 0)
	var headRes *cos.Response
	if access == types.CsAccessPrivate {
		headRes, err = tc.HeadWithPrivate(s3Key)
	} else {
		headRes, err = tc.HeadWithPublic(s3Key)
	}

	if err == nil && headRes.StatusCode == http.StatusOK {
		// 同一文件反复上传
		logs.Info("[uploadCore] duplicate upload of the same file, s3key: %s, err: %v", s3Key, err)
		code = cerror.CodeSuccess
		return
	} else {
		logs.Info("[uploadCore] need upload file, s3key: %s, err: %v", s3Key, err)
	}

	if access == types.CsAccessPrivate {
		err = tc.UploadFromStream2Private(s3Key, f)
	} else {
		err = tc.UploadFromStream2Public(s3Key, f)
	}

	if err != nil {
		logs.Error("[uploadCore] upload to cos fail. file:", upFilename, ", err:", err)
		code = cerror.UploadResourceFail
		return
	}

	code = cerror.CodeSuccess

	return
}

func (c *WebAdmBaseController) UploadResource(upFilename string, access int) (s3Key string, code cerror.ErrCode, mediaInfo types.MediaSimpleInfo, err error) {
	code = cerror.CodeSuccess

	if access != types.CsAccessPrivate && access < types.CsAccessPublic {
		code = cerror.AccessDenied
		logs.Warn("[UploadResource] access denied, can't upload file. for:", upFilename, "err:", err)
		return
	}

	f, h, err := c.GetFile(upFilename)
	if err != nil {
		code = cerror.PermissionDenied
		logs.Warn("[UploadResource] permission denied, can't upload file. for:", upFilename, "err:", err)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	s3Key, _, code, mediaInfo, err = c.uploadCore(f, h, upFilename, access)
	return s3Key, code, mediaInfo, err
}

// GetString 简单封装,防止xss攻击
func (c *WebAdmBaseController) GetString(key string, def ...string) string {
	into := c.Controller.GetString(key, def...)
	into = libtools.RegRemoveScript(into)
	into = strings.TrimSpace(into)
	return into
}
