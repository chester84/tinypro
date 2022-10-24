package controllers

import (
	"tinypro/common/cerror"

	"github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	web.Controller
}

func (c *ErrorController) Prepare() {
}

func (c *ErrorController) Error404() {
	c.Data["json"] = cerror.BuildApiResponse(cerror.ApiNotFound, cerror.EmptyData)
	c.ServeJSON()
}

func (c *ErrorController) Error501() {
	c.Data["json"] = cerror.BuildApiResponse(cerror.ServiceUnavailable, cerror.EmptyData)
	c.ServeJSON()
}

func (c *ErrorController) Error500() {
	c.Data["json"] = cerror.BuildApiResponse(cerror.ServiceUnavailable, cerror.EmptyData)
	c.ServeJSON()
}

func (c *ErrorController) ErrorDb() {
	c.Data["json"] = cerror.BuildApiResponse(cerror.ServiceUnavailable, cerror.EmptyData)
	c.ServeJSON()
}
