package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
	"tinypro/common/cerror"
	"tinypro/common/types"
)

type envArgsT struct {
	Api   string
	Env   string
	Token string

	Rid        string
	AppSN      int
	PaySN      string
	OrderSN    string
	UserSN     string
	Url        string
	Amount     string
	ObjSN      string
	OpenUserID string
	Mobile     string
	SmsCode    string
	Code       string
	Page       int
	Size       int
}

var envArgs envArgsT

func init() {
	flag.StringVar(&envArgs.Api, "api", "",
		`API URL
# TOOLS
/web-api/address2geo    transfer address to geo
/web-api/geo2address    transfer geo to address

/web-api/ping
/web-api/oauth-login    the party login
/web-api/get-config     get configuration


/web-api/op-broadcast          
/web-api/op-feed               


/web-api/front-page            

/open-api/wx-js-config 
/web-api/mnp-check-version
`)
	flag.StringVar(&envArgs.Env, "env", "dev", "dev,test,prod")
	flag.StringVar(&envArgs.Token, "token", "", "token")
	flag.IntVar(&envArgs.AppSN, "app-sn", 1, "setup app sn, default 0")
	flag.StringVar(&envArgs.OrderSN, "order-sn", "", "ordersn")
	flag.StringVar(&envArgs.PaySN, "pay-sn", "", "paysn")
	flag.StringVar(&envArgs.UserSN, "user-sn", "", "userid")
	flag.StringVar(&envArgs.OpenUserID, "open-user-id", "", "openid")
	flag.StringVar(&envArgs.ObjSN, "obj-sn", "", "objsn")
	flag.StringVar(&envArgs.Mobile, "mobile", "", "mobile")
	flag.StringVar(&envArgs.Amount, "amount", "10", "amount")
	flag.StringVar(&envArgs.SmsCode, "sms-code", "", "sms-code")
	flag.StringVar(&envArgs.Rid, "rid", "", "rid")
	flag.StringVar(&envArgs.Code, "code", "", "code")
	flag.IntVar(&envArgs.Page, "page", 1, "page")
	flag.IntVar(&envArgs.Size, "size", 10, "size")
	flag.StringVar(&envArgs.Url, "url", "", "url")
}

func main() {
	flag.Parse()

	logs.Debug("debug web-api ...")

	var host string
	switch envArgs.Env {
	case "dev":
		host = "http://127.0.0.1:8976"

	case "test":
		host = ""

	case "prod":
		host = ""

	default:
		flag.PrintDefaults()
		logs.Error("env is wrong")
		os.Exit(0)
	}

	request := map[string]interface{}{
		"trace_id":     libtools.Md5(libtools.GenerateRandomStr(16)),
		"request_time": fmt.Sprintf(`%d`, libtools.GetUnixMillis()),
		"access_token": envArgs.Token,
		"app_sn":       envArgs.AppSN,
	}

	var apiUrl string
	switch envArgs.Api {

	case "/web-api/ping":
		request["page"] = 1
		request["size"] = 10

	case "/web-api/address2geo":
		request["address"] = ``

	case "/web-api/geo2address":
		request["latitude"] = `39.984154`
		request["longitude"] = `116.307490`

	case "/web-api/mnp-check-version":
		request["app_version"] = `1.2.3`

	case "/web-api/oauth-login":
		request["open_oauth_plt"] = types.OpenOauthWeChat
		request["open_user_id"] = "of1Vg5xKvNGlxNqFw3SOSnJf0ii0"
		request["wx_open_id"] = `xxooyy`
		request["nickname"] = "wx-002"
		request["nickname"] = "test"
		//request["open_avatar"] = `https://dcydb3a2bq85h.cloudfront.net/uc/cover/default0.png`

	default:
		flag.PrintDefaults()
		logs.Error("please specify the api")
		os.Exit(0)
	}

	apiUrl = host + envArgs.Api
	logs.Notice(">>> ENV: %s, API: %s <<<", envArgs.Env, apiUrl)

	reqHeaders := map[string]string{
		"Connection":   "keep-alive",
		"Content-Type": "application/json",
		"X-Trace":      types.XTrace,
		"User-Agent":   "go-api/v1 app/1.0.0.004/4",
		//"User-Agent": "com.luanchen.NearbyMoments/1.0.0.1001/ (iPhone; iOS 13.1.1; Scale/3.00)",
		//"User-Agent": "okhttp/4.3.1 (Android 10) Covid-19 Map/1.0.0.004/4",
		//"User-Agent": "okhttp/4.3.1 (Android 9) Panditji/1.0.1.32/42",
	}

	var reqData string

	reqJSON, _ := libtools.JsonEncode(request)
	reqData = reqJSON

	logs.Info("[trace] request data: %v", reqData)
	httpBody, httpStatusCode, err := libtools.SimpleHttpClient("POST", apiUrl, reqHeaders, reqData, libtools.DefaultHttpTimeout())
	logs.Notice("httpBody: %s, httpStatusCode: %d, err: %v", httpBody, httpStatusCode, err)

	var apiData cerror.ApiResponse
	err = json.Unmarshal(httpBody, &apiData)
	if apiData.Code == cerror.CodeSuccess {
		logs.Alert("api ok")
	} else {
		logs.Critical("api data[error code: %d] err", apiData.Code)
	}

	logs.Info("debug process has done.")
}
