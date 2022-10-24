package routers

import (
	"tinypro/api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/web-api/wx/login-register", &controllers.WebApiSiWeiController{}, "post:WxOauthLoginOrRegister")
	beego.Router("/web-api/config", &controllers.WebApiSiWeiController{}, "post:Config")
	beego.Router("/web-api/front-page", &controllers.WebApiSiWeiController{}, "post:FrontPage")
	beego.Router("/web-api/public-courses", &controllers.WebApiSiWeiController{}, "post:PublicCourses")
	beego.Router("/web-api/last-enroll", &controllers.WebApiSiWeiController{}, "post:GetLastEnroll")
	beego.Router("/web-api/enroll", &controllers.WebApiSiWeiController{}, "post:Enroll")
	beego.Router("/web-api/enroll/is-pass", &controllers.WebApiSiWeiController{}, "post:IsPassEnroll")
	beego.Router("/web-api/enroll/my-enrolls", &controllers.WebApiSiWeiController{}, "post:MyEnrolls")
	beego.Router("/web-api/enroll/sign-in", &controllers.WebApiSiWeiController{}, "post:SignInCourse")
	beego.Router("/web-api/enroll/my-courses", &controllers.WebApiSiWeiController{}, "post:MyCourses")
	beego.Router("/web-api/enroll/my-history-courses", &controllers.WebApiSiWeiController{}, "post:MyHistoryCourses")
	beego.Router("/web-api/enroll/detail", &controllers.WebApiSiWeiController{}, "post:GetEnrollDetail")
	beego.Router("/web-api/my-msg-list", &controllers.WebApiSiWeiController{}, "post:MyMsgList")
	beego.Router("/web-api/satis-survey", &controllers.WebApiSiWeiController{}, "post:SatisSurvey")
	beego.Router("/web-api/subscribe-template", &controllers.WebApiSiWeiController{}, "post:SubscribeTemplate")
	beego.Router("/web-api/advertise-popup", &controllers.WebApiSiWeiController{}, "post:AdvertisePopup")
}
