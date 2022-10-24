package controllers

import (
	beego "github.com/beego/beego/v2/server/web"

	"tinypro/common/cerror"
	"github.com/chester84/libtools"
	"tinypro/common/types"
)

type MainController struct {
	beego.Controller
}

// 以下的两个路由,不走统一的加解密,参数签名检查

func (c *MainController) Get() {
	res := cerror.ApiResponse{
		Code:      cerror.CodeSuccess,
		Message:   "What are you doing?",
		SeverTime: libtools.GetUnixMillis(),
		Data: struct {
		}{},
	}

	c.Data["json"] = res
	c.ServeJSON()
}

func (c *MainController) Ping() {
	data := map[string]interface{}{
		"server_time": libtools.GetUnixMillis(),
		"version":     types.AppVersion,
		"head_hash":   libtools.GitRevParseHead(),
		//"router":      c.Ctx.Request.RequestURI,
		//"ip": c.Ctx.Input.IP(),
	}
	res := cerror.ApiResponse{
		Code:      cerror.CodeSuccess,
		Message:   "pong",
		Data:      data,
		SeverTime: libtools.GetUnixMillis(),
	}

	c.Data["json"] = res
	c.ServeJSON()
}

func (c *MainController) Favicon() {
	_ = c.Ctx.Output.Body([]byte(``))
}

func (c *MainController) MPVerify() {
	_ = c.Ctx.Output.Body([]byte(``))
}
