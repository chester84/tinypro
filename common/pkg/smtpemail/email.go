package smtpemail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

const (
	smtpHost = "smtp.163.com"
	smtpPort = 465
	smtpUser = "cli_mail@163.com"
	smtpPwd  = "YEOHQCXIDWYXRUVV"
)

func DefaultReceiver() []string {
	return []string{
		"kefu@huimer.com",
		"liuliang@huimer.com",
	}
}

func Send(fromAlias string, to []string, subject string, body string, contentType string) (err error) {
	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPwd,
		smtpHost,
	)

	message := composeMsg(fromAlias, to, subject, body, contentType)

	err = sendMailUsingTLS(
		fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth,
		smtpUser,
		to,
		[]byte(message),
	)

	if err != nil {
		logs.Error("send get exception, err: %v", err)
	}

	return
}

// compose message according to "from, to, subject, body"
func composeMsg(fromAlias string, to []string, subject string, body string, contentType string) (message string) {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = fromAlias + "<" + smtpUser + ">"
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	if contentType == "html" {
		headers["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}
	// Setup message
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return
}

// return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		logs.Error("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 参考net/smtp的func SendMail()
// 使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
// len(to)>1时,to[1]开始提示是密送
func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	//create smtp client
	c, err := dial(addr)
	if err != nil {
		logs.Error("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logs.Error("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
