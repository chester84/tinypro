package routers

import (
	"tinypro/api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// /api/v1 group
	beego.Router("/api/v1/ping", &controllers.MainController{}, "*:Ping")
	beego.Router("/api/v1/encrypt/ping", &controllers.APIV1Controller{}, "post:Ping")

	// web-api group
	beego.Router("/web-api/mnp-check-version", &controllers.WebApiToolsController{}, "post:MnpCheckVersion")
	//// 工具
	beego.Router("/web-api/address2geo", &controllers.WebApiToolsController{}, "post:Address2Geo")
	beego.Router("/web-api/geo2address", &controllers.WebApiToolsController{}, "post:Geo2Address")
	beego.Router("/web-api/ping", &controllers.WebApiController{}, "post:Ping")
	beego.Router("/web-api/oauth-login", &controllers.WebApiController{}, "post:OauthLogin")
	beego.Router("/web-api/op-broadcast", &controllers.WebApiController{}, "post:OpBroadcast")
	beego.Router("/web-api/update-profile", &controllers.WebApiController{}, "post:UpdateProfile")
	//// 常见问题
	beego.Router("/web-api/faq", &controllers.WebApiController{}, "post:Faq")
	beego.Router("/web-api/faq/more", &controllers.WebApiController{}, "post:FaqMore")
	beego.Router("/web-api/get-config", &controllers.WebApiController{}, "post:GetConfig")
	// 微信
	beego.Router("/web-api/wx/decrypt", &controllers.WebApiController{}, "post:WxDecrypt")
}
