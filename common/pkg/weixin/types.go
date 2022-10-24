package weixin

import (
	"github.com/chester84/libtools"

	"github.com/beego/beego/v2/core/logs"
)

const (
	AppSNWxGzh = 0 // 微信公众号
	AppSNWxMng = 1 // 微信小程序
)

func AppID(appSN int) string {
	var appid string
	switch appSN {
	case AppSNWxGzh:
		appid = gzhAppID()

	case AppSNWxMng:
		appid = MnpAppID()

	default:
		logs.Error("[AppID] get unexpected appSN: %d", appSN)
	}

	return appid
}

func Secret(appSN int) string {
	var secret string

	switch appSN {
	case AppSNWxGzh:
		secret = gzhSecret()

	case AppSNWxMng:
		secret = MnpSecret()

	default:
		logs.Error("[Secret] get unexpected appSN: %d", appSN)
	}

	return secret
}

func gzhAppID() string {
	if libtools.IsProductEnv() {
		return ``
	} else {
		return ``
	}
}

func gzhSecret() string {
	if libtools.IsProductEnv() {
		return ``
	} else {
		return ``
	}
}

// 四维小程序appId
func MnpAppID() string {
	return `wx44701b8b29557073`
}

// 四维小程序secret
func MnpSecret() string {
	return `7394c7cce7c185428930333efc1b680a`
}

/*
商户号
*/
func MerchantID() string {
	return ""
}

func MerchantProdKey() string {
	return "" // 保护码:
}

/*
PayKey
*/
func MerchantKey() string {
	if libtools.IsProductEnv() {
		return MerchantProdKey()
	} else {
		return GetSandBoxPayMerchantKey()
	}
}

type AppMessageShareData struct {
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Link   string `json:"link"`
	ImgUrl string `json:"img_url"`
}
