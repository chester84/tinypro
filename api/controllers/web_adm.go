package controllers

import (
	"bytes"

	_ "image/png"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/captcha"
	"github.com/teris-io/shortid"

	"tinypro/common/cerror"
	"tinypro/common/models"
	"tinypro/common/pkg/accesstoken"
	"tinypro/common/pkg/admin"
	"tinypro/common/pkg/rbacv2"
	"tinypro/common/pkg/tc"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type WebAdmController struct {
	WebAdmBaseController
}

func (c *WebAdmController) Prepare() {
	// 调用上一级的 Prepare 方
	c.WebAdmBaseController.Prepare()
}

func (c *WebAdmController) Ping() {
	data := map[string]interface{}{
		"server_time": libtools.GetUnixMillis(),
		"version":     types.AdminVersion,
		"head_hash":   libtools.GitRevParseHead(),
		"router":      c.Ctx.Request.RequestURI,
		"ip":          c.IP,
		"user_agent":  c.Ctx.Input.UserAgent(),
	}
	c.SuccessResponse(data)
}

func (c *WebAdmController) LoginCaptcha() {
	if c.AccountId > 0 {
		c.Ctx.Output.Status = 403
		_ = c.Ctx.Output.Body([]byte(types.Html403))
		return
	}

	captchaStr := libtools.Int2Str(libtools.GenerateRandom(100000, 999999))
	cookieValue, err := shortid.Generate()
	if err != nil {
		logs.Error("[LoginCaptcha] generate short id get exception, err: %v", err)
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	// 写 cookie
	c.Ctx.SetCookie(types.AdmLoginCaptchaCookieName, cookieValue, 0, "/")
	// 写头
	c.Ctx.Output.Header("Content-Type", "image/png")

	digits := make([]byte, 6)
	for i, _ := range captchaStr {
		digits[i] = captchaStr[i] - 48
	}

	captchaImg := captcha.NewImage(digits, 240, 80)
	buf := new(bytes.Buffer)
	_, err = captchaImg.WriteTo(buf)
	if err != nil {
		logs.Error("[LoginCaptcha] write img to bug get exception, err: %v", err)
		c.Ctx.Output.Status = 404
		_ = c.Ctx.Output.Body([]byte(types.Html404))
		return
	}

	// 写缓存
	admin.SetLoginCaptcha(cookieValue, captchaStr)

	_ = c.Ctx.Output.Body(buf.Bytes())
}

func (c *WebAdmController) Login() {
	userName := c.GetString(`nickname`)
	password := c.GetString(`password`)

	adm, err := admin.OneByEmail(userName)
	if err != nil {
		logs.Error("[Login] wrong input, ip: %s, userName: %s", c.IP, userName)
		c.TerminateWithCode(cerror.AccessDenied)
		return
	}

	check := admin.CheckLoginIsValid(userName, password)
	if !check {
		logs.Error("[Login] wrong input, ip: %s, userName: %s", c.IP, userName)
		c.TerminateWithCode(cerror.InvalidMobileOrPassword)
		return
	}

	accessToken, err := accesstoken.GenTokenWithCache(adm.Id, types.PlatformAdm, c.IP)
	if err != nil {
		logs.Error("[Login] gen token exception, userId: %d, ip: %s, err: %v", adm.Id, c.IP, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	type resT struct {
		AccessToken string `json:"access_token"`
	}

	res := resT{
		AccessToken: accessToken,
	}

	c.SuccessResponse(res)
}

func (c *WebAdmController) Logout() {
	accessToken := c.Ctx.Request.Header.Get("X-Access-Token")
	accesstoken.CleanTokenCache(types.PlatformAdm, accessToken)

	c.SuccessResponse(types.H{
		"op_msg": `操作成功`,
	})
}

func (c *WebAdmController) BaseInfo() {
	roleNameBox, _ := rbacv2.OperatorAccessRole(c.AccountId)

	type resT struct {
		Nickname     string   `json:"nickname"`
		Avatar       string   `json:"avatar"`
		Roles        []string `json:"roles"`
		Introduction string   `json:"introduction"`
	}

	res := resT{
		Nickname: c.Account.Nickname,
		Avatar:   `https://cdn-1302993108.cos.ap-guangzhou.myqcloud.com/img/default.png`,
		Roles:    roleNameBox,
	}
	c.SuccessResponse(res)
}

func (c *WebAdmController) ChangePassword() {
	originPassword := c.GetString(`originPassword`)
	newPassword := c.GetString(`newPassword`)
	if len(originPassword) < 6 || len(newPassword) < 6 || originPassword == newPassword {
		c.TerminateWithCode(cerror.InvalidPassword)
		return
	}

	originE := libtools.PasswordEncrypt(originPassword, c.Account.RegisterTime)
	if originE != c.Account.Password {
		logs.Warning("[ChangePassword] 原始密码错误. adminUid: %d", c.AccountId)
		c.TerminateWithCode(cerror.InvalidOldPassword)
		return
	}

	newE := libtools.PasswordEncrypt(newPassword, c.Account.RegisterTime)
	c.Account.Password = newE
	_, err := models.OrmUpdate(&c.Account, []string{`Password`})
	if err != nil {
		logs.Error("[ChangePassword] db update exception, admin: %#v, err: %v", c.Account, err)
		c.TerminateWithCode(cerror.ServiceDbOpFail)
		return
	}

	c.SuccessResponse("")
}

func (c *WebAdmController) UploadResource4Tinymce() {
	s3key, code, mediaInfo, err := c.UploadResource("file", types.CsAccessPublic)
	if err != nil {
		c.TerminateWithCode(code)
		return
	}

	resourceUrl := tc.PublicUrl(s3key)

	type resT struct {
		ResourceUrl string `json:"resource_url"`
		types.MediaSimpleInfo
	}

	res := resT{
		ResourceUrl:     resourceUrl,
		MediaSimpleInfo: mediaInfo,
	}

	c.SuccessResponse(res)
}
