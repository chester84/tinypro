package main

import (
	"bytes"
	"fmt"
	"github.com/chester84/libtools"
	"time"
	_ "tinypro/common/lib/db/mysql"
	"tinypro/common/pkg/tc"

	"github.com/beego/beego/v2/core/logs"
	goqrcode "github.com/skip2/go-qrcode"
)

func main() {
	qrcodeCreate()
}

func qrcodeCreate() {
	q, err := goqrcode.New(fmt.Sprintf("https://www.baidu.com?t=%d", time.Now().Unix()), goqrcode.Medium)
	if err != nil {
		return
	}

	png, err := q.PNG(150)
	if err != nil {
		return
	}

	fileMd5 := libtools.Md5Bytes(png)
	_, s3Key := libtools.BuildHashName(fileMd5, "png")

	imgBuf := new(bytes.Buffer)
	imgBuf.Write(png)

	err = tc.UploadFromStream2Public(s3Key, imgBuf)
	if err != nil {
		logs.Warn("[convertImageUrl] UploadByFileByte err: %v", err)
	}

	url := tc.PublicUrl(s3Key)
	logs.Debug("url is %s", url)
}
