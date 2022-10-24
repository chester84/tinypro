package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"tinypro/api/controllers"
	"tinypro/common/cprof"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "*:Get")
	beego.Router("/ping", &controllers.MainController{}, "*:Ping")
	beego.Router("/favicon.ico", &controllers.MainController{}, "*:Favicon")

	// 微信
	beego.Router("/MP_VERIFY", &controllers.MainController{}, "*:MPVerify")
	beego.Router("/wx/entrust", &controllers.WXNotifyController{}, "*:EntrustNotify")
	//beego.Router("/wx/unified-order/notify", &controllers.WXNotifyController{}, "*:UnifiedOrderPayNotify")

	// pprof 性能分析
	beego.Router("/debug/pprof", &cprof.ProfController{}, "*:Get")
	beego.Router(`/debug/pprof/:pp([\w]+)`, &cprof.ProfController{}, "*:Get")

	// open-api
	beego.Router("/open-api/wx-oauth2-silent", &controllers.OpenApiController{}, "*:WxOauth2Silent")
	beego.Router("/open-api/resource/:rid", &controllers.ResourceController{}, "*:Resource")
	beego.Router("/open-api/wx/msg/callback", &controllers.WXNotifyController{}, "*:MsgCallback")
	// 暂时放开授权限制
	beego.Router("/open-api/wx-js-config", &controllers.OpenApiController{}, "post:WxJsConfig")

	beego.Router("/ws/v1", &controllers.WebsocketController{}, "get:Get")

}
