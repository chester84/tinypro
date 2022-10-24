package main

import (
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// 数据库初始化
	_ "tinypro/common/lib/db/mysql"
	// 路由
	_ "tinypro/api/routers"

	"github.com/chester84/libtools"
	"tinypro/api/controllers"
	"tinypro/common/cerror"
	"tinypro/common/lib/clogs"
)

var FilterRouter = func(ctx *context.Context) {
	logs.Info("step into filter. reqs:%s", ctx.Request.URL.Path)

	ip := ctx.Input.IP()
	if !libtools.IsInternalIPV1(ip) {
		resObj := cerror.BuildApiResponse(cerror.AccessDenied, cerror.EmptyData)
		resJSON, _ := json.Marshal(resObj)
		ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
		_, _ = ctx.ResponseWriter.Write(resJSON)
		return
	}
}

func main() {
	dir, _ := config.String("log_dir")
	port, _ := config.String("httpport")
	clogs.InitLog(dir, "api_"+port)

	logs.Info("start api.")

	web.Handler("/metrics", promhttp.Handler())

	web.ErrorController(&controllers.ErrorController{})

	var success = []byte("SUPPORT OPTIONS")

	var corsFunc = func(ctx *context.Context) {
		origin := ctx.Input.Header("Origin")
		ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
		ctx.Output.Header("Access-Control-Max-Age", "3600")
		ctx.Output.Header("Access-Control-Allow-Headers", "X-Custom-Header,accept,Content-Type,Access-Token,X-Access-Token")
		ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		ctx.Output.Header("Access-Control-Allow-Origin", origin)
		if ctx.Input.Method() == http.MethodOptions {
			ctx.Output.SetStatus(http.StatusOK)
			_ = ctx.Output.Body(success)
		}
	}

	web.InsertFilter("/*", web.BeforeRouter, corsFunc)

	web.InsertFilter("/metrics", web.BeforeRouter, FilterRouter)

	web.Run()
}
