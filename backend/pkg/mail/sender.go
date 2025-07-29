package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/hasan-kilici/chat/pkg/utils"
)

var config utils.ConfigList

func init() {
	config = utils.LoadConfig("./configs/service.ini")
}

var (
	password  = config.SMTPPassword
	from      = config.SMTPFrom
	SMTPHost  = config.SMTPHost
	SMTPPort  = config.SMTPPort
	GmailAuth = smtp.PlainAuth("", from, password, SMTPHost)
)

func SendMailTemplate(to string, subject string, templatePath string, data any) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("template parse hatası: %w", err)
	}

	var body bytes.Buffer

	headers := "MIME-version: 1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, headers)))

	err = t.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	err = smtp.SendMail(SMTPHost+":"+SMTPPort, GmailAuth, from, []string{to}, body.Bytes())
	if err != nil {
		return fmt.Errorf("mail send error: %w", err)
	}

	return nil
}
