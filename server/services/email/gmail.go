package email

import (
	"os"

	"gopkg.in/gomail.v2"
	"schej.it/server/logger"
)

func isGmailConfigured() bool {
	return os.Getenv("GMAIL_APP_PASSWORD") != "" && os.Getenv("SCHEJ_EMAIL_ADDRESS") != ""
}

type gmailProvider struct{}

func (p *gmailProvider) sendTemplate(to string, tmpl Template, data map[string]any) error {
	subject, err := render(tmpl.Subject, data)
	if err != nil {
		return err
	}
	body, err := render(tmpl.Body, data)
	if err != nil {
		return err
	}

	fromEmail := os.Getenv("SCHEJ_EMAIL_ADDRESS")
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")

	m := gomail.NewMessage()
	if tmpl.FromEmail != "" {
		m.SetHeader("From", tmpl.FromEmail)
	} else {
		m.SetHeader("From", fromEmail)
	}
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, appPassword)
	if err := d.DialAndSend(m); err != nil {
		logger.StdErr.Println(err)
		return err
	}

	return nil
}

func (p *gmailProvider) ensureSubscriber(_ string, _ bool)               {}
func (p *gmailProvider) subscriberExists(_ string) (bool, *int)          { return false, nil }
func (p *gmailProvider) addSubscriber(_, _, _, _ string, _ *int, _ bool) {}
