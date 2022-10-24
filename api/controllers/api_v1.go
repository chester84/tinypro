package controllers

import (
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type APIV1Controller struct {
	APIBaseController
}

func (c *APIV1Controller) Prepare() {
	// 调用上一级的 Prepare 方
	c.APIBaseController.Prepare()
}

func (c *APIV1Controller) Ping() {
	// 测试表明,库表报缺字段的错误时, beego 给的 http 状态码是 200,之前注册的 500 拦截器不会生效
	//abc := c.RequestJSON["abc"].(string)
	//logs.Debug(abc)

	data := map[string]interface{}{
		"server_time": libtools.GetUnixMillis(),
		"version":     types.AppVersion,
		"head_hash":   libtools.GitRevParseHead(),
		"router":      c.Ctx.Request.RequestURI,
		"ip":          c.Ctx.Input.IP(),
		"user_agent":  c.Ctx.Input.UserAgent(),
	}
	c.SuccessResp(data)
}
