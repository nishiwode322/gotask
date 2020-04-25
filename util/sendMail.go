package util

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailParam struct {
	ServerHost string

	ServerPort int

	FromEmail string

	FromPasswd string

	Toers string

	CCers string
}

var serverHost, fromEmail, fromPasswd string
var serverPort int

func InitEmail(ep *EmailParam, m *gomail.Message) {

	toers := []string{}
	if len(ep.Toers) == 0 {
		return
	}

	for _, tmp := range strings.Split(ep.Toers, ",") {
		toers = append(toers, strings.TrimSpace(tmp))
	}

	m.SetHeader("To", toers...)

	if len(ep.CCers) != 0 {
		for _, tmp := range strings.Split(ep.CCers, ",") {
			toers = append(toers, strings.TrimSpace(tmp))
		}
		m.SetHeader("Cc", toers...)
	}

	m.SetAddressHeader("From", ep.FromEmail, "")
}

func SendEmail(subject, body string, ep *EmailParam, m *gomail.Message) {

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer(ep.ServerHost, ep.ServerPort, ep.FromEmail, ep.FromPasswd)

	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func RunSendMail() {

	serverHost := "smtp.exmail.qq.com"
	serverPort := 465
	fromEmail := "cicd@latelee.org"
	fromPasswd := "1qaz@WSX"

	myToers := "li@latelee.org, latelee@163.com"
	myCCers := ""

	//set mail context
	subject := "这是主题"
	body := `这是正文<br>
            <h3>这是标题</h3>
             Hello <a href = "http://www.latelee.org">主页</a><br>`
	//
	myEmail := &EmailParam{
		ServerHost: serverHost,
		ServerPort: serverPort,
		FromEmail:  fromEmail,
		FromPasswd: fromPasswd,
		Toers:      myToers,
		CCers:      myCCers,
	}
	m := gomail.NewMessage()
	InitEmail(myEmail, m)
	SendEmail(subject, body, myEmail, m)
}
