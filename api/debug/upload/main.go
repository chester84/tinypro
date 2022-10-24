package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beego/beego/v2/core/logs"

	"github.com/chester84/libtools"
	"tinypro/common/cerror"
)

type envArgsT struct {
	Api   string
	Env   string
	Token string
	File  string
}

var envArgs envArgsT

// 本地
const dftDevAccessToken = `9c7e1b690d215861ed709008bf9ce5cc`

// 线上
const dftProdAccessToken = `da079e3c183fb7606166a36c289094cd`

var accessToken string

func init() {
	flag.StringVar(&envArgs.Api, "api", "",
		`api url
/api/v1/upload/resource
`)
	flag.StringVar(&envArgs.Env, "env", "dev", "dev,test,prod")
	flag.StringVar(&envArgs.Token, "token", "", "token")
	flag.StringVar(&envArgs.File, "file", "./debug/upload/home.png", "upload files")

}

func main() {
	flag.Parse()

	logs.Debug("debug api ...")

	var host string
	switch envArgs.Env {
	case "dev":
		host = "http://127.0.0.1:8565"
		accessToken = dftDevAccessToken

	case "test":
		host = ""

	case "prod":
		host = ""
		accessToken = dftProdAccessToken

	default:
		flag.PrintDefaults()
		logs.Error("env is wrong")
		os.Exit(0)
	}

	if envArgs.Token != "" {
		accessToken = envArgs.Token
	}

	buf := new(bytes.Buffer)
	var etag string

	var apiUrl string
	switch envArgs.Api {

	case "/api/v1/upload/resource":
		content, err := ioutil.ReadFile(envArgs.File)
		if err != nil {
			logs.Error("can not read file: %s, err: %v", envArgs.File, err)
			os.Exit(0)
		}
		buf.Write(content)
		etag = libtools.Md5Bytes(content)

	default:
		flag.PrintDefaults()
		logs.Error("please specify the api")
		os.Exit(0)
	}
	apiUrl = host + envArgs.Api

	logs.Notice("-----API: %s\n", apiUrl)

	reqHeaders := map[string]string{
		"Connection":     "keep-alive",
		"User-Agent":     "go-api/v1",
		"X-Access-Token": accessToken,
		"X-Request-Time": fmt.Sprintf(`%d`, libtools.GetUnixMillis()),
		"X-ETag":         etag,
	}

	logs.Debug("HttpHeader: %v", reqHeaders)

	httpBody, httpStatusCode, err := libtools.SimpleHttpClient("POST", apiUrl, reqHeaders, buf.String(), libtools.DefaultHttpTimeout())
	logs.Notice("httpBody: %s, httpStatusCode: %d, err: %v\n", httpBody, httpStatusCode, err)

	var apiData cerror.ApiResponse
	err = json.Unmarshal(httpBody, &apiData)
	if apiData.Code == cerror.CodeSuccess {
		logs.Notice("apiResData: %v", apiData.Data)
	} else {
		logs.Notice("api[error code]error.\n")
	}
}
