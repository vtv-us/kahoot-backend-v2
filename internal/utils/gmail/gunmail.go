package gmail

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type MailgunService struct {
	Client    *mailgun.MailgunImpl
	EmailFrom string
	Host      string
}

func NewGunMailService(c *utils.Config) MailgunService {
	client := mailgun.NewMailgun(c.GunmailAddress, c.GunmailApiKey)
	return MailgunService{
		Client:    client,
		EmailFrom: c.Mail,
		Host:      c.ServerAddress,
	}
}

func (s *MailgunService) SendEmailForVerified(toEmail string, verifyString string) error {
	sender := s.EmailFrom
	subject := "Verify your Kahoot account"
	body := fmt.Sprintf(`Click on the following link to verify your account: %s/auth/verify/%s/%s`, s.Host, toEmail, verifyString)
	recipient := toEmail

	msg := s.Client.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err := s.Client.Send(ctx, msg)

	return err
}
