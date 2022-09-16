package email

import (
	"testing"
	"time"

	"gopkg.in/gomail.v2"
)

func TestSend(t *testing.T) {
	sender := "your@email.com" // 发件人
	port := 465                // 发件端口
	authcode := "authcode"     // 授权码/密码
	c := New(
		"smtp.qq.com", // 发件箱
		port,
		sender,
		authcode,
	)
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {sender},
		"To":      {sender},
		"Subject": {"测试邮件"},
	})
	m.SetBody("text/html", "<h1>测试内容5555</h1>")
	c.Send(m)
	<-time.After(5 * time.Second)
}
