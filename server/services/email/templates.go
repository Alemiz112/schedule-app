package email

import (
	"os"
	"strconv"
)

// Template defines an email for any provider.
// Subject and Body are Go text/template strings; map keys are accessed as {{.keyName}}.
type Template struct {
	// listmonkIDEnv is the env var that holds the Listmonk template ID for this email.
	listmonkIDEnv string
	// listmonkIDDefault is the fallback Listmonk template ID when the env var is unset.
	listmonkIDDefault int
	// FromEmail optionally overrides the sender address (Listmonk + Gmail).
	FromEmail string
	// Subject is a Go text/template string used by non-Listmonk providers.
	Subject string
	// Body is a Go text/template string (plain text) used by non-Listmonk providers.
	Body string
}

func (t Template) listmonkTemplateID() int {
	if t.listmonkIDEnv != "" {
		if v := os.Getenv(t.listmonkIDEnv); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return t.listmonkIDDefault
}

// Templates holds every outbound email the application sends.
// Listmonk template IDs are read from env vars at send time (with hardcoded defaults
// matching the values that were previously inline in the route handlers).
var Templates = struct {
	OTP                   Template
	GroupInvite           Template
	AddedAttendee         Template
	SomeoneRespondedEvent Template
	SomeoneRespondedGroup Template
	XResponses            Template
	EveryoneResponded     Template
	AppointmentRequested  Template
	AppointmentApproved   Template
}{
	OTP: Template{
		listmonkIDEnv: "LISTMONK_OTP_EMAIL_TEMPLATE_ID",
		FromEmail:     "Timeful <noreply@timeful.app>",
		Subject:       "Your Timeful verification code",
		Body:          "Your verification code is: {{.code}}\n\nThis code expires in 10 minutes.",
	},
	GroupInvite: Template{
		listmonkIDEnv:     "LISTMONK_GROUP_INVITE_TEMPLATE_ID",
		listmonkIDDefault: 9,
		Subject:           "{{.ownerName}} invited you to {{.groupName}}",
		Body:              "{{.ownerName}} invited you to fill out \"{{.groupName}}\".\n\n{{.groupUrl}}",
	},
	AddedAttendee: Template{
		listmonkIDEnv:     "LISTMONK_ADDED_ATTENDEE_TEMPLATE_ID",
		listmonkIDDefault: 11,
		Subject:           "New people added to {{.groupName}}",
		Body:              "{{.ownerName}} added new people to \"{{.groupName}}\".\n\n{{.groupUrl}}",
	},
	SomeoneRespondedEvent: Template{
		listmonkIDEnv:     "LISTMONK_SOMEONE_RESPONDED_EVENT_TEMPLATE_ID",
		listmonkIDDefault: 10,
		Subject:           "{{.respondentName}} responded to {{.eventName}}",
		Body:              "{{.respondentName}} just responded to your event \"{{.eventName}}\".\n\n{{.eventUrl}}",
	},
	SomeoneRespondedGroup: Template{
		listmonkIDEnv:     "LISTMONK_SOMEONE_RESPONDED_GROUP_TEMPLATE_ID",
		listmonkIDDefault: 13,
		Subject:           "{{.respondentName}} responded to {{.groupName}}",
		Body:              "{{.respondentName}} just responded to your group \"{{.groupName}}\".\n\n{{.groupUrl}}",
	},
	XResponses: Template{
		listmonkIDEnv:     "LISTMONK_X_RESPONSES_TEMPLATE_ID",
		listmonkIDDefault: 14,
		Subject:           "{{.numResponses}} people responded to {{.eventName}}",
		Body:              "{{.numResponses}} people have responded to your event \"{{.eventName}}\".\n\n{{.eventUrl}}",
	},
	EveryoneResponded: Template{
		listmonkIDEnv:     "LISTMONK_EVERYONE_RESPONDED_TEMPLATE_ID",
		listmonkIDDefault: 8,
		Subject:           "Everyone responded to {{.eventName}}",
		Body:              "Everyone has responded to your event \"{{.eventName}}\".\n\n{{.eventUrl}}",
	},
	AppointmentRequested: Template{
		listmonkIDEnv: "LISTMONK_APPOINTMENT_REQUESTED_TEMPLATE_ID",
		Subject:       "New appointment request for {{.eventName}}",
		Body:          "{{.requesterName}} has requested an appointment for \"{{.eventName}}\" on {{.date}} from {{.startTime}} to {{.endTime}}.\n\nManage requests: {{.eventUrl}}",
	},
	AppointmentApproved: Template{
		listmonkIDEnv: "LISTMONK_APPOINTMENT_APPROVED_TEMPLATE_ID",
		Subject:       "Your appointment has been confirmed",
		Body:          "Hi {{.requesterName}},\n\nYour appointment for \"{{.eventName}}\" on {{.date}} from {{.startTime}} to {{.endTime}} has been confirmed.\n\n{{.eventUrl}}",
	},
}
