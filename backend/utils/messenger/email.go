package messenger

import (
	"backend/global"
	_ "html/template"

	"gopkg.in/gomail.v2"
)

// Sent message for customer
func SendEmail(recipients []string, subject string, message string) (err error) {

	m := gomail.NewMessage()
	m.SetHeader("From", global.Conf.SmtpUser)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	//m.SetBody("text/html", template.HTMLEscapeString(Message))
	m.SetBody("text/html", message)
	d := gomail.NewDialer(global.Conf.SmtpHost, global.Conf.SmtpPort, global.Conf.SmtpUser, global.Conf.SmtpPassword)

	if err := d.DialAndSend(m); err != nil {
		//log.Error("Failed to send email ", err)
		return err
	}
	return
}
