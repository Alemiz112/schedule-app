package email

import (
	"bytes"
	"sync"
	"text/template"

	"schej.it/server/logger"
)

// Provider is the interface for email delivery backends.
type Provider interface {
	sendTemplate(to string, tmpl Template, data map[string]any) error
	ensureSubscriber(email string, sendMarketing bool)
	subscriberExists(email string) (bool, *int)
	addSubscriber(email, firstName, lastName, picture string, subscriberID *int, sendMarketing bool)
}

var (
	initOnce sync.Once
	provider Provider
)

func getProvider() Provider {
	initOnce.Do(func() {
		switch {
		case isListmonkConfigured():
			provider = &listmonkProvider{}
		case isGmailConfigured():
			provider = &gmailProvider{}
		default:
			provider = &noopProvider{}
		}
	})
	return provider
}

// Send sends a templated email to the given address.
func Send(to string, tmpl Template, data map[string]any) {
	if err := getProvider().sendTemplate(to, tmpl, data); err != nil {
		logger.StdErr.Println("email:", err)
	}
}

// SendAndAddSubscriber sends a templated email, first ensuring the recipient
// is a subscriber (Listmonk-specific; a no-op for other providers).
func SendAndAddSubscriber(to string, tmpl Template, data map[string]any, sendMarketing bool) {
	p := getProvider()
	p.ensureSubscriber(to, sendMarketing)
	if err := p.sendTemplate(to, tmpl, data); err != nil {
		logger.StdErr.Println("email:", err)
	}
}

// SubscriberExists reports whether the given email is a subscriber (Listmonk only).
func SubscriberExists(email string) (bool, *int) {
	return getProvider().subscriberExists(email)
}

// AddSubscriber adds a subscriber to the mailing list (Listmonk only, no-op for others).
func AddSubscriber(email, firstName, lastName, picture string, subscriberID *int, sendMarketing bool) {
	getProvider().addSubscriber(email, firstName, lastName, picture, subscriberID, sendMarketing)
}

// render applies a Go text/template string to a data map.
func render(tmplStr string, data map[string]any) (string, error) {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// noopProvider silently discards all emails.
type noopProvider struct{}

func (p *noopProvider) sendTemplate(_ string, _ Template, _ map[string]any) error { return nil }
func (p *noopProvider) ensureSubscriber(_ string, _ bool)                         {}
func (p *noopProvider) subscriberExists(_ string) (bool, *int)                    { return false, nil }
func (p *noopProvider) addSubscriber(_, _, _, _ string, _ *int, _ bool)           {}
