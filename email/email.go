package email

import (
	"fmt"
	"net/http"

	"github.com/Abhi-singh-karuna/my_Liberary/baselogger"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Email interface defines the structure for an email.
type Email interface {
	Address() string
	Subject() *string
	HTML() *string
	Text() *string
}

// RawEmail interface defines the structure for a raw email.
type RawEmail interface {
	Body() []byte
}

// email struct implements the Email interface.
type email struct {
	address string
	subject *string
	html    *string
	text    *string
}

func (e *email) Address() string {
	return e.address
}

func (e *email) Subject() *string {
	return e.subject
}

func (e *email) HTML() *string {
	return e.html
}

func (e *email) Text() *string {
	return e.text
}

func NewEmail(a string, s, h, t *string) Email {
	return &email{a, s, h, t}
}

// rawEmail struct implements the RawEmail interface.
type rawEmail struct {
	body []byte
}

func (e *rawEmail) Body() []byte {
	return e.body
}

func NewRawEmail(body []byte) RawEmail {
	return &rawEmail{body}
}

// Sender interface defines the methods for sending emails.
type Sender interface {
	Send(Email) error
	SendRaw(RawEmail) error
}

// SendGridSender struct implements the Sender interface using SendGrid.
type EmailService struct {
	apiKey       string
	fromEmail    string
	fromName     string
	MailRequired bool
	logger       *baselogger.BaseLogger
}

// apiKey: cfg.SendGridAPIKey, fromEmail: cfg.SendGridFromEmail, fromName:
func SendGridEmailService(apiKey string, fromEmail string, fromName string, MailRequired bool, logger *baselogger.BaseLogger) *EmailService {
	return &EmailService{apiKey: apiKey, fromEmail: fromEmail, fromName: fromName, MailRequired: MailRequired, logger: logger}
}

// Send sends a structured email using SendGrid.
func (s *EmailService) Send(e Email) error {
	if s.MailRequired {
		from := mail.NewEmail(s.fromName, s.fromEmail)
		to := mail.NewEmail("", e.Address())
		message := mail.NewSingleEmail(from, *e.Subject(), to, *e.Text(), *e.HTML())
		client := sendgrid.NewSendClient(s.apiKey)
		response, err := client.Send(message)
		if err != nil || response.StatusCode != http.StatusAccepted {
			return fmt.Errorf("failed to send email: %v, response code: %d", err, response.StatusCode)
		}
	}
	return nil
}

// SendRaw sends a raw email using SendGrid.
func (s *EmailService) SendRaw(e RawEmail) error {
	request := sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = e.Body()
	response, err := sendgrid.API(request)
	if err != nil || response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send raw email: %v, response code: %d", err, response.StatusCode)
	}
	return nil
}
