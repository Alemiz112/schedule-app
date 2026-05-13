package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/services/calendar"
	emailsvc "schej.it/server/services/email"
	"schej.it/server/utils"
)

// @Summary Submit an appointment request for a time slot
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{startDate=string,endDate=string,name=string,email=string,notes=string} true "Appointment request details"
// @Success 201 {object} models.AppointmentRequest
// @Router /events/{eventId}/appointment-requests [post]
func createAppointmentRequest(c *gin.Context) {
	payload := struct {
		StartDate primitive.DateTime `json:"startDate" binding:"required"`
		EndDate   primitive.DateTime `json:"endDate" binding:"required"`
		Name      string             `json:"name" binding:"required"`
		Email     string             `json:"email" binding:"required"`
		Notes     string             `json:"notes"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	if !utils.Coalesce(event.IsAppointment) {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotAppointment})
		return
	}

	// Enforce max appointments limit
	if event.MaxAppointments != nil && *event.MaxAppointments > 0 {
		if db.GetActiveAppointmentCount(event.Id.Hex()) >= *event.MaxAppointments {
			c.JSON(http.StatusConflict, responses.Error{Error: errs.AppointmentLimitReached})
			return
		}
	}

	req := models.AppointmentRequest{
		Id:        primitive.NewObjectID(),
		EventId:   event.Id,
		StartDate: payload.StartDate,
		EndDate:   payload.EndDate,
		Name:      payload.Name,
		Email:     payload.Email,
		Notes:     payload.Notes,
		Status:    models.AppointmentPending,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := db.AppointmentRequestsCollection.InsertOne(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Error: "insert-failed"})
		return
	}

	// Cancel any pending reminder emails for this booker and persist
	if req.Email != "" {
		cancelRemindersForEmail(event, req.Email)
		db.EventsCollection.UpdateByID(context.Background(), event.Id, bson.M{"$set": event})
	}

	// Notify owner of new request
	owner := db.GetUserById(event.OwnerId.Hex())
	if owner != nil && owner.Email != "" {
		eventUrl := fmt.Sprintf("%s/e/%s", utils.GetBaseUrl(), event.GetId())
		emailsvc.Send(owner.Email, emailsvc.Templates.AppointmentRequested, appointmentEmailData(&req, event, eventUrl))
	}

	// Auto-approve if enabled and owner has calendar auto-add configured
	if utils.Coalesce(event.AutoApproveAppointments) {
		if owner != nil && owner.CalendarOptions != nil && owner.CalendarOptions.AddToCalendar {
			db.UpdateAppointmentRequestStatus(req.Id.Hex(), models.AppointmentApproved)
			req.Status = models.AppointmentApproved

			eventUrl := fmt.Sprintf("%s/e/%s", utils.GetBaseUrl(), event.GetId())

			if req.Email != "" {
				emailsvc.Send(req.Email, emailsvc.Templates.AppointmentApproved, appointmentEmailData(&req, event, eventUrl))

			}

			calendarKey := owner.CalendarOptions.DefaultCalendarKey
			if calendarKey == "" && owner.PrimaryAccountKey != nil {
				calendarKey = *owner.PrimaryAccountKey
			}
			if calendarKey != "" {
				attendeeEmails := []string{}
				if req.Email != "" {
					attendeeEmails = append(attendeeEmails, req.Email)
				}
				calendarId := owner.CalendarOptions.DefaultCalendarId
				calendar.CreateEventForUser(owner, calendarKey, calendarId, calendar.CreateCalendarEventInput{
					Title:          fmt.Sprintf("%s with %s", event.Name, req.Name),
					StartDate:      req.StartDate.Time(),
					EndDate:        req.EndDate.Time(),
					Description:    fmt.Sprintf("Booked via Timeful: %s", eventUrl),
					AttendeeEmails: attendeeEmails,
				})
			}
		}
	}

	c.JSON(http.StatusCreated, req)
}

// @Summary Get all appointment requests for an event (owner only)
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param status query string false "Filter by status (pending, approved, rejected)"
// @Success 200 {array} models.AppointmentRequest
// @Router /events/{eventId}/appointment-requests [get]
func getAppointmentRequests(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, responses.Error{Error: errs.UserNotEventOwner})
		return
	}

	if !utils.Coalesce(event.IsAppointment) {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotAppointment})
		return
	}

	var statusFilter *models.AppointmentStatus
	if statusParam := c.Query("status"); statusParam != "" {
		s := models.AppointmentStatus(statusParam)
		statusFilter = &s
	}

	reqs := db.GetAppointmentRequestsByEvent(event.Id.Hex(), statusFilter)
	c.JSON(http.StatusOK, reqs)
}

// @Summary Approve an appointment request (owner only)
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param requestId path string true "Appointment request ID"
// @Success 200 {object} models.AppointmentRequest
// @Router /events/{eventId}/appointment-requests/{requestId}/approve [post]
func approveAppointmentRequest(c *gin.Context) {
	req, event := getAppointmentRequestWithOwnerCheck(c)
	if req == nil || event == nil {
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	db.UpdateAppointmentRequestStatus(req.Id.Hex(), models.AppointmentApproved)
	req.Status = models.AppointmentApproved

	eventUrl := fmt.Sprintf("%s/e/%s", utils.GetBaseUrl(), event.GetId())
	if req.Email != "" {
		emailsvc.Send(req.Email, emailsvc.Templates.AppointmentApproved, appointmentEmailData(req, event, eventUrl))
	}

	calendarEventCreated := false
	if user.CalendarOptions != nil && user.CalendarOptions.AddToCalendar {
		calendarKey := user.CalendarOptions.DefaultCalendarKey
		if calendarKey == "" && user.PrimaryAccountKey != nil {
			calendarKey = *user.PrimaryAccountKey
		}
		if calendarKey != "" {
			attendeeEmails := []string{}
			if req.Email != "" {
				attendeeEmails = append(attendeeEmails, req.Email)
			}
			calendarId := user.CalendarOptions.DefaultCalendarId
			err := calendar.CreateEventForUser(user, calendarKey, calendarId, calendar.CreateCalendarEventInput{
				Title:          fmt.Sprintf("%s with %s", event.Name, req.Name),
				StartDate:      req.StartDate.Time(),
				EndDate:        req.EndDate.Time(),
				Description:    fmt.Sprintf("Booked via Timeful: %s", eventUrl),
				AttendeeEmails: attendeeEmails,
			})
			if err == nil {
				calendarEventCreated = true
			}
		}
	}

	type approveResponse struct {
		models.AppointmentRequest
		CalendarEventCreated bool `json:"calendarEventCreated"`
	}
	c.JSON(http.StatusOK, approveResponse{
		AppointmentRequest:   *req,
		CalendarEventCreated: calendarEventCreated,
	})
}

// @Summary Reject an appointment request (owner only)
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param requestId path string true "Appointment request ID"
// @Success 200
// @Router /events/{eventId}/appointment-requests/{requestId}/reject [post]
func rejectAppointmentRequest(c *gin.Context) {
	req, event := getAppointmentRequestWithOwnerCheck(c)
	if req == nil || event == nil {
		return
	}

	db.UpdateAppointmentRequestStatus(req.Id.Hex(), models.AppointmentRejected)
	req.Status = models.AppointmentRejected

	c.JSON(http.StatusOK, req)
}

// @Summary Get booked time slots for an appointment event (public — no personal info)
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} object{startDate=string,endDate=string}
// @Router /events/{eventId}/appointment-requests/booked [get]
func getBookedAppointmentSlots(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	if !utils.Coalesce(event.IsAppointment) {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotAppointment})
		return
	}

	type BookedSlot struct {
		StartDate primitive.DateTime `json:"startDate"`
		EndDate   primitive.DateTime `json:"endDate"`
	}

	reqs := db.GetAppointmentRequestsByEvent(event.Id.Hex(), nil)
	slots := make([]BookedSlot, 0)
	for _, req := range reqs {
		if req.Status == models.AppointmentPending || req.Status == models.AppointmentApproved {
			slots = append(slots, BookedSlot{StartDate: req.StartDate, EndDate: req.EndDate})
		}
	}

	c.JSON(http.StatusOK, slots)
}

// shared validation: fetch request, verify event ownership, verify request belongs to event
func getAppointmentRequestWithOwnerCheck(c *gin.Context) (*models.AppointmentRequest, *models.Event) {
	eventId := c.Param("eventId")
	requestId := c.Param("requestId")

	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return nil, nil
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, responses.Error{Error: errs.UserNotEventOwner})
		return nil, nil
	}

	req := db.GetAppointmentRequest(requestId)
	if req == nil || req.EventId != event.Id {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.AppointmentRequestNotFound})
		return nil, nil
	}

	return req, event
}

// @Summary Get the event owner's available time slots from their event response (public)
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} string
// @Router /events/{eventId}/appointment-requests/owner-availability [get]
func getOwnerAvailability(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	if !utils.Coalesce(event.IsAppointment) {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotAppointment})
		return
	}

	// Find the owner's response among all event responses
	ownerIdHex := event.OwnerId.Hex()
	for _, er := range db.GetEventResponses(event.Id.Hex()) {
		if er.UserId == ownerIdHex && er.Response != nil {
			// Return availability slots (nil response means owner hasn't responded yet)
			slots := er.Response.Availability
			if slots == nil {
				slots = []primitive.DateTime{}
			}
			c.JSON(http.StatusOK, slots)
			return
		}
	}

	// Owner hasn't submitted availability — return nil so the frontend shows all slots open
	c.JSON(http.StatusOK, nil)
}

func appointmentEmailData(req *models.AppointmentRequest, event *models.Event, eventUrl string) map[string]any {
	loc := time.Now().Location()
	start := req.StartDate.Time().In(loc)
	end := req.EndDate.Time().In(loc)
	return map[string]any{
		"requesterName": req.Name,
		"eventName":     event.Name,
		"date":          start.Format("January 2, 2006"),
		"startTime":     start.Format("3:04 PM"),
		"endTime":       end.Format("3:04 PM"),
		"eventUrl":      eventUrl,
		"notes":         req.Notes,
	}
}
