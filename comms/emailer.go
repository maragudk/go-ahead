package comms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-ahead/errors2"
)

// Email is an unparsed email address.
type Email = string

// NameAndEmail combo, of the form "Name <email@example.com>"
type NameAndEmail = string

const (
	defaultEndpoint = "https://api.postmarkapp.com"
)

// Emailer can send emails through Postmark.
type Emailer struct {
	token    string
	endpoint string
}

func NewEmailer(token string) *Emailer {
	return &Emailer{
		token:    token,
		endpoint: defaultEndpoint,
	}
}

type sendWithTemplateRequestBody struct {
	From          NameAndEmail
	To            NameAndEmail
	TemplateAlias string
	TemplateModel map[string]string
}

// SendWithTemplate using the Postmark API with a named template and template data.
func (e *Emailer) SendWithTemplate(ctx context.Context, from, to NameAndEmail, templateAlias string, templateModel map[string]string) error {
	body := sendWithTemplateRequestBody{
		From:          from,
		To:            to,
		TemplateAlias: templateAlias,
		TemplateModel: templateModel,
	}
	return e.sendRequest(ctx, "/email/withTemplate", body, e.token)
}

type sendRequestBody struct {
	From     NameAndEmail
	To       NameAndEmail
	Subject  string
	HtmlBody string
}

// Send using the Postmark API.
func (e *Emailer) Send(ctx context.Context, from, to NameAndEmail, subject, htmlBody string) error {
	body := sendRequestBody{
		From:     from,
		To:       to,
		Subject:  subject,
		HtmlBody: htmlBody,
	}
	return e.sendRequest(ctx, "/email", body, e.token)
}

func (e *Emailer) sendRequest(ctx context.Context, endpoint string, body interface{}, token string) error {
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return errors2.Wrap(err, "could not marshal request body to json")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint+endpoint, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return errors2.Wrap(err, "could not build email request")
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Postmark-Server-Token", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors2.Wrap(err, "could not do email request")
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode > 299 {
		return fmt.Errorf("could not send email, got status %v", response.StatusCode)
	}

	return nil
}

// CreateNameAndEmail returns a name and email string ready for inserting into From and To fields.
func CreateNameAndEmail(name string, email Email) NameAndEmail {
	return fmt.Sprintf("%v <%v>", name, email)
}
