package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppointmentStatus string

const (
	AppointmentPending  AppointmentStatus = "pending"
	AppointmentApproved AppointmentStatus = "approved"
	AppointmentRejected AppointmentStatus = "rejected"
)

type AppointmentRequest struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	EventId primitive.ObjectID `json:"eventId" bson:"eventId,omitempty"`

	StartDate primitive.DateTime `json:"startDate" bson:"startDate,omitempty"`
	EndDate   primitive.DateTime `json:"endDate" bson:"endDate,omitempty"`

	// Guest or logged-in user info
	Name   string             `json:"name" bson:"name,omitempty"`
	Email  string             `json:"email" bson:"email,omitempty"`
	Notes  string             `json:"notes" bson:"notes,omitempty"`
	UserId primitive.ObjectID `json:"userId" bson:"userId,omitempty"`

	Status    AppointmentStatus  `json:"status" bson:"status"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt,omitempty"`
}
