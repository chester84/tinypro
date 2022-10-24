package ckeditor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"

	"tinypro/common/pkg/system/config"
	"tinypro/common/pkg/tc"
	"github.com/chester84/libtools"
)

// PreHandleRichAttrImg 处理服务文本
// 1 图片宽度不能大于800
// 2 外链图片转为本地
func PreHandleRichAttrImg(content string) (string, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(content))
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return content, fmt.Errorf("NewDocumentFromReader err: %v", err)
	}

	doc.Find("img").Each(ImgHandler)
	doc.Find("style").Each(RemoveHandler)
	doc.Find("link").Each(RemoveHandler)
	doc.Find("meta").Each(RemoveHandler)

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return content, fmt.Errorf("OuterHtml err: %v", err)
	}

	return out, nil
}

func ImgHandler(index int, this *goquery.Selection) {
	// 1 宽度处理
	style, exist := this.Attr("style")
	if exist {
		logs.Info("[imgHandler] got img style: %v", style)

		// style := "height:949px; width:2331px"
		var h, w int
		_, _ = fmt.Sscanf(style, "height: %dpx; width: %dpx", &h, &w)

		// stupid magic number width:800
		magic := 800
		if w > magic {
			style = strings.Replace(style, fmt.Sprintf("height: %dpx;", h), "", -1)
			style = strings.Replace(style, fmt.Sprintf("width: %dpx", w), fmt.Sprintf("width: %dpx;", magic), -1)
			this.SetAttr("style", style)
		}
	}

	// 2 外链处理
	src, exist := this.Attr("src")
	if exist {
		logs.Info("[imgHandler] got img src: %v", src)
		// 图片地址不是自己cdn 且图片路径维权路径的 外链转换
		if !strings.HasPrefix(src, tc.CDNDomain()) &&
			(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
			dst := convertImageUrl(src)
			logs.Info("src: %s, after convert url is: %v", src, dst)
			this.SetAttr("src", dst)
		}
	}
}

func convertImageUrl(src string) string {
	img, code, err := libtools.SimpleHttpClient(libtools.HttpMethodGet, src, libtools.DefaultReqHeaders(), "", libtools.DefaultHttpTimeout())
	if err != nil || code != 200 {
		logs.Warn("[convertImageUrl] src: %v SimpleHttpClient err: %v", src, err)
		return src
	}

	if len(img) < 300 {
		logs.Error("[convertImageUrl] img Len: %v, src: %v", len(img), src)
		return src
	}

	fileMd5 := libtools.Md5Bytes(img)
	_, s3Key := libtools.BuildHashName(fileMd5, "png")
	imgBuf := new(bytes.Buffer)
	imgBuf.Write(img)
	err = tc.UploadFromStream2Public(s3Key, imgBuf)
	if err != nil {
		logs.Warn("[convertImageUrl] UploadByFileByte src: %v, err: %v", src, err)
		return src
	}

	rs := tc.PublicUrl(s3Key)
	return rs
}

func RemoveHandler(index int, this *goquery.Selection) {
	this.Remove()
}

// InsertCss 将CSS加到头部
func insertCss(content, css string) string {
	if content == "" {
		logs.Warning("[insertCss] content is empty")
		return ""
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte(content))
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return content
	}

	// todo 判断富文本是否有内容
	body := doc.Find("body")
	if body == nil || isHtmlEmpty(body.Clone()) {
		return ""
	} else {
		b, _ := body.Html()
		logs.Notice("[insertCss] b: %v", b)
		body.SetHtml("<div>" + b + "</div>")
		script := config.ValidItemString("_rich_text_internal_script")
		if script != "" {
			body.AfterHtml(script)
		}
	}

	head := doc.Find("head")
	if head != nil {
		head.AppendHtml("<meta charset='UTF-8'>")
		head.AppendHtml(css)

	} else {
		return css + content
	}

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return content
	}
	return out
}

// InsertCss 将CSS加到头部
func OnlyStyleLink(content string) (string, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte(content))
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return content, fmt.Errorf("css style err: %v", err)
	}

	r1 := doc.Find("style").Map(only)
	r2 := doc.Find("link").Map(only)
	r3 := doc.Find("meta").Map(only)

	r1 = append(r1, r2...)
	r1 = append(r3, r1...)

	return strings.Join(r1, ""), nil
}

func only(index int, this *goquery.Selection) string {
	out, err := goquery.OuterHtml(this)
	if err != nil {
		return ""
	}
	return out
}

func isHtmlEmpty(body *goquery.Selection) bool {
	if body == nil {
		return true
	}

	b := body
	b.Find("script").Remove()

	t := b.Find("table").Remove()
	if len(t.Nodes) > 0 {
		logs.Info("[isHtmlEmpty] find table so no empty.")
		return false
	}

	i := b.Find("img").Remove()
	if len(i.Nodes) > 0 {
		logs.Info("[isHtmlEmpty] find img so no empty.")
		return false
	}

	f := b.Find("iframe").Remove()
	if len(f.Nodes) > 0 {
		logs.Info("[isHtmlEmpty] find iframe so no empty.")
		return false
	}

	ps := b.Find("pre").Map(getText)
	v := strings.TrimSpace(strings.Join(ps, ""))
	if v != "" {
		logs.Info("[isHtmlEmpty] find pre so no empty.")
		return false
	}

	rs := b.Find("p").Map(getText)
	logs.Notice(libtools.JsonEncode(rs))
	v = strings.TrimSpace(strings.Join(rs, ""))
	logs.Info("[isHtmlEmpty] find getText v: %v", v)

	return v == ""
}

func getText(index int, this *goquery.Selection) string {
	v := this.Text()
	// 看看v里是不是有干货，还是只有回车和换行
	vv := strings.Replace(v, "\n", "", -1)
	vv = strings.Replace(vv, "\t", "", -1)
	vv = strings.TrimSpace(vv)
	if len(vv) == 0 {
		return ""
	}
	return v
}
