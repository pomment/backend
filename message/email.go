package message

import (
	"context"
	"fmt"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/config"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/k3a/html2text"
	"github.com/mailgun/mailgun-go/v4"
)

func SendEmailMailgun(sendEmailContext common.SendEmailContext) (err error) {
	emailConfig := config.Content.Email
	mg := mailgun.NewMailgun(emailConfig.MailgunDomain, emailConfig.MailgunAPIKey)

	sender := emailConfig.Sender
	subject := sendEmailContext.Subject
	body := sendEmailContext.Body
	recipient := sendEmailContext.Recipient

	message := mg.NewMessage(sender, subject, html2text.HTML2Text(body), recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	fmt.Printf("Email sent via Mailgun API. ID: %s Resp: %s\n", id, resp)
	return nil
}

func SendEmailSMTP(sendEmailContext common.SendEmailContext) (err error) {
	emailConfig := config.Content.Email
	e := email.NewEmail()
	e.From = config.Content.Email.Sender
	e.To = []string{sendEmailContext.Recipient}
	e.Subject = sendEmailContext.Subject
	e.Text = []byte(html2text.HTML2Text(sendEmailContext.Body))
	e.HTML = []byte(sendEmailContext.Body)
	err = e.Send(fmt.Sprintf("%s:%d", emailConfig.SMTPHost, emailConfig.SMTPPort), smtp.PlainAuth("", emailConfig.SMTPUsername, emailConfig.SMTPPassword, emailConfig.SMTPHost))
	return err
}

func SendEmail(to string, post common.Post, parentPost common.Post, thread common.Thread) (err error) {
	emailConfig := config.Content.Email
	template := config.Content.WebTemplate
	fmt.Printf("Attempting to send email to %s", to)

	if !emailConfig.Enabled {
		return nil
	}

	emailTemplateConfig := map[string]interface{}{
		"Post":           post,
		"ParentPost":     parentPost,
		"Thread":         thread,
		"UnsubscribeURL": "https://www.tcdw.net",
	}
	emailContent := common.SendEmailContext{
		Recipient: to,
		Subject:   template.EmailTitle.Render(emailTemplateConfig),
		Body:      template.EmailBody.Render(emailTemplateConfig),
	}
	switch emailConfig.Mode {
	case "mailgun":
		err = SendEmailMailgun(emailContent)
		break
	case "smtp":
		err = SendEmailSMTP(emailContent)
		break
	}
	return err
}
