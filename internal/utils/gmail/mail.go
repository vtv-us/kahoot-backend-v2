package gmail

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type SendgridService struct {
	Client    *sendgrid.Client
	EmailFrom string
	Host      string
	Frontend  string
}

type EmailContent struct {
	From             *mail.Email
	To               *mail.Email
	Subject          string
	PlainTextContent string
	HtmlContent      string
}

func NewMailService(c *utils.Config) SendgridService {
	client := sendgrid.NewSendClient(c.SendgridApiKey)
	return SendgridService{
		Client:    client,
		EmailFrom: c.SendgridEmail,
		Host:      c.ServerAddress,
		Frontend:  c.FrontendAddress,
	}
}

func (s *SendgridService) SendEmailForVerified(toEmail string, verifyString string) error {
	emailContent := EmailContent{
		From: &mail.Email{
			Name:    "Kahoot",
			Address: s.EmailFrom,
		},
		To: &mail.Email{
			Name:    "User",
			Address: toEmail,
		},
		Subject:          "Verify your Kahoot account",
		PlainTextContent: fmt.Sprintf(`Click on the following link to verify your account: %s/auth/verify/%s/%s`, s.Host, toEmail, verifyString),
		HtmlContent:      fmt.Sprintf(`<p>Click on the following link to verify your account: <a href="%s/auth/verify/%s/%s">link</a></p>`, s.Host, toEmail, verifyString),
	}
	message := mail.NewSingleEmail(emailContent.From, emailContent.Subject, emailContent.To, emailContent.PlainTextContent, emailContent.HtmlContent)
	_, err := s.Client.Send(message)
	return err
}

func (s *SendgridService) SendEmailForInvite(toEmail string, groupID string, groupName string, inviter string) error {
	emailContent := EmailContent{
		From: &mail.Email{
			Name:    "Kahoot",
			Address: s.EmailFrom,
		},
		To: &mail.Email{
			Name:    "User",
			Address: toEmail,
		},
		Subject:          "You have been invited to a Kahoot group",
		PlainTextContent: fmt.Sprintf(`You have been invited to join the group %s by %s. Click on the following link to join: %s`, groupName, inviter, utils.GenLink(s.Frontend, groupID)),
		HtmlContent:      fmt.Sprintf(`<p>You have been invited to join the group %s by %s. Click on the following link to join: <a href="%s">link</a></p>`, groupName, inviter, utils.GenLink(s.Frontend, groupID)),
	}
	message := mail.NewSingleEmail(emailContent.From, emailContent.Subject, emailContent.To, emailContent.PlainTextContent, emailContent.HtmlContent)
	_, err := s.Client.Send(message)
	return err
}
