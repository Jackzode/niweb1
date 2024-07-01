package handler

import (
	"gopkg.in/gomail.v2"
	"log"
	"testing"
)

func TestEmailService_Send(t *testing.T) {

	from_email := "to.be.migrant@gmail.com"
	//from_name := "auto.migrant"
	//smtp_host := "smtp.gmail.com"
	//smtp_port := 465
	smtp_password := "kqtb mhfv bfuw pdop"
	//smtp_username := "to.be.migrant@gmail.com"
	//encryption := "SSL"

	// 设置收件人信息
	to := "jackzhi942@gmail.com"

	// 设置邮件内容
	subject := "Test Email - 465"
	body := "This is a test email sent from Golang using gomail."

	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", from_email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 配置SMTP服务器
	d := gomail.NewDialer("smtp.gmail.com", 465, from_email, smtp_password)
	//d.SSL = true
	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
