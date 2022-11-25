package gmail

import (
	"crypto/tls"
	"fmt"

	"github.com/vtv-us/kahoot-backend/internal/utils"
	gomail "gopkg.in/gomail.v2"
)

type GomailService struct {
	Client    *gomail.Dialer
	EmailFrom string
	Host      string
}

func NewMailServiceV2(c *utils.Config) GomailService {
	client := gomail.NewDialer("smtp.gmail.com", 587, c.Mail, c.MailPassword)
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return GomailService{
		Client:    client,
		EmailFrom: c.Mail,
		Host:      c.ServerAddress,
	}
}

func (s *GomailService) SendEmailForVerified(toEmail string, verifyString string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.EmailFrom)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Verify your Kahoot account")
	msg.SetBody("text/html", fmt.Sprintf(`<p>Click on the following link to verify your account: <a href="%s/auth/verify/%s/%s">link</a></p>`, s.Host, toEmail, verifyString))

	return s.Client.DialAndSend(msg)
}
