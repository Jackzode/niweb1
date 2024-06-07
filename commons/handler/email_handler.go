package handler

import (
	"crypto/tls"
	"github.com/Jackzode/painting/commons/constants"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/Jackzode/painting/config"
	"gopkg.in/gomail.v2"
	"os"
)

var EmailHandler *EmailService

type EmailService struct {
	g *gomail.Dialer
}

func InitEmailService(config *config.EmailConfig) {
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	if config.Encryption == constants.SSL {
		d.SSL = true
	}
	if len(os.Getenv("SKIP_SMTP_TLS_VERIFY")) > 0 {
		d.TLSConfig = &tls.Config{ServerName: d.Host, InsecureSkipVerify: true}
	}
	EmailHandler = &EmailService{g: d}
	return
}

// Send email send
func (es *EmailService) Send(toEmailAddr, title, body string) {

	m := gomail.NewMessage()
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	err := es.g.DialAndSend(m)
	if err != nil {
		glog.Slog.Infof("send email to %s success", toEmailAddr)
	}
}
