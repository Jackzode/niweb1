package handler

import (
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/config"
	"gopkg.in/gomail.v2"
)

var EmailHandler *EmailService

type EmailService struct {
	from string
	g    *gomail.Dialer
}

func InitEmailService(config *config.EmailConfig) {

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	if config.Encryption == constants.SSL {
		d.SSL = true
	}
	//if len(os.Getenv("SKIP_SMTP_TLS_VERIFY")) > 0 {
	//	d.TLSConfig = &tls.Config{ServerName: d.Host, InsecureSkipVerify: true}
	//}
	EmailHandler = &EmailService{g: d}
	EmailHandler.from = config.FromEmail
	return
}

// Send email send
func (es *EmailService) Send(toEmailAddr, title, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	m.SetHeader("To", toEmailAddr)
	m.SetHeader("From", es.from)
	err := es.g.DialAndSend(m)
	return err
}
