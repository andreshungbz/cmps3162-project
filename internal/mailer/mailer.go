package mailer

import (
	"bytes"
	"embed"
	"time"

	"github.com/wneessen/go-mail"

	ht "html/template"
	tt "text/template"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer defines the client connection to the SMTP server and sender information.
type Mailer struct {
	client *mail.Client
	sender string
}

// New creates a new instance of Mailer with configurations set.
func New(host string, port int, username, password, sender string) (*Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second), // timeout for sending email
	)
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		client: client,
		sender: sender,
	}

	return mailer, nil
}

// Send dynamically sets a template to send as an email to a recipient. It sends
// both a plaintext and HTML version.
func (m *Mailer) Send(recipient string, templateFile string, data any) error {
	// Template Construction

	// parse text/template
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// execute subject template with dynamic data
	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// execute plainBody template with dynamic data
	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// parse html/template
	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// execute htmlBody template with dynamic data
	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Message Construction

	msg := mail.NewMsg()

	// set recipient email
	err = msg.To(recipient)
	if err != nil {
		return err
	}

	// set sender email
	err = msg.From(m.sender)
	if err != nil {
		return nil
	}

	// set email subject
	msg.Subject(subject.String())
	// set plaintext body string
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	// set HTML body string as an alternative
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	// make 3 attempts to send an email in case of network issues
	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}

		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}
