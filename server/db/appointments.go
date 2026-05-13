package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/logger"
	"schej.it/server/models"
)

func GetAppointmentRequest(requestId string) *models.AppointmentRequest {
	objectId, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		return nil
	}

	result := AppointmentRequestsCollection.FindOne(context.Background(), bson.M{"_id": objectId})
	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}

	var req models.AppointmentRequest
	if err := result.Decode(&req); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &req
}

func GetAppointmentRequestsByEvent(eventId string, statusFilter *models.AppointmentStatus) []models.AppointmentRequest {
	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return []models.AppointmentRequest{}
	}

	filter := bson.M{"eventId": objectId}
	if statusFilter != nil {
		filter["status"] = *statusFilter
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := AppointmentRequestsCollection.Find(context.Background(), filter, opts)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	reqs := make([]models.AppointmentRequest, 0)
	if err := cursor.All(context.Background(), &reqs); err != nil {
		logger.StdErr.Panicln(err)
	}

	return reqs
}

// GetActiveAppointmentCount returns the number of pending+approved requests for an event.
func GetActiveAppointmentCount(eventId string) int {
	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		return 0
	}
	count, err := AppointmentRequestsCollection.CountDocuments(context.Background(), bson.M{
		"eventId": objectId,
		"status":  bson.M{"$in": bson.A{models.AppointmentPending, models.AppointmentApproved}},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return int(count)
}

func UpdateAppointmentRequestStatus(requestId string, status models.AppointmentStatus) {
	objectId, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		return
	}

	_, err = AppointmentRequestsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
}
