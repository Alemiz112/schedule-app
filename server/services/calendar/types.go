package calendar

import (
	"time"

	"schej.it/server/models"
)

type CreateCalendarEventInput struct {
	Title          string
	StartDate      time.Time
	EndDate        time.Time
	Description    string
	AttendeeEmails []string
	CalendarId     string // sub-calendar ID; providers fall back to primary/default if empty
}

type CalendarProvider interface {
	GetCalendarList() (map[string]models.SubCalendar, error)
	GetCalendarEvents(calendarId string, timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, error)
	CreateCalendarEvent(input CreateCalendarEventInput) error
}

func GetCalendarProvider(calendarAccount models.CalendarAccount) CalendarProvider {
	switch calendarAccount.CalendarType {
	case models.GoogleCalendarType:
		return &GoogleCalendar{
			OAuth2CalendarAuth: *calendarAccount.OAuth2CalendarAuth,
		}
	case models.OutlookCalendarType:
		return &OutlookCalendar{
			OAuth2CalendarAuth: *calendarAccount.OAuth2CalendarAuth,
		}
	case models.AppleCalendarType:
		return &AppleCalendar{
			AppleCalendarAuth: *calendarAccount.AppleCalendarAuth,
		}
	case models.ICSCalendarType:
		return &ICSCalendar{
			ICSCalendarAuth: *calendarAccount.ICSCalendarAuth,
		}
	}
	return nil
}
