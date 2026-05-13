package scheduler

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/logger"
	emailsvc "schej.it/server/services/email"
	"schej.it/server/utils"
)

type scheduledEmail struct {
	Id          primitive.ObjectID `bson:"_id"`
	ScheduledAt time.Time          `bson:"scheduledAt"`
	To          string             `bson:"to"`
	TemplateKey string             `bson:"templateKey"`
	Data        bson.M             `bson:"data"`
}

// templateRegistry maps a stable string key to the live Template value.
// Keys are persisted in MongoDB so never rename them.
var templateRegistry = map[string]emailsvc.Template{
	"ReminderInitial":            emailsvc.Templates.ReminderInitial,
	"ReminderSecond":             emailsvc.Templates.ReminderSecond,
	"ReminderFinal":              emailsvc.Templates.ReminderFinal,
	"AppointmentReminderInitial": emailsvc.Templates.AppointmentReminderInitial,
	"AppointmentReminderSecond":  emailsvc.Templates.AppointmentReminderSecond,
	"AppointmentReminderFinal":   emailsvc.Templates.AppointmentReminderFinal,
}

// Start launches the background goroutine that polls for due emails every minute.
func Start() {
	go func() {
		for {
			sendDue()
			time.Sleep(time.Minute)
		}
	}()
}

// CreateEmailReminders sends the initial reminder immediately and schedules
// the 24 h and 72 h follow-ups. Returns the IDs of the two scheduled jobs.
func CreateEmailReminders(email, ownerName, eventName, eventId string, isAppointment bool) []string {
	baseUrl := utils.GetBaseUrl()

	var initialTmpl, secondKey, finalKey string
	var data bson.M

	if isAppointment {
		initialTmpl = "AppointmentReminderInitial"
		secondKey = "AppointmentReminderSecond"
		finalKey = "AppointmentReminderFinal"
		data = bson.M{
			"ownerName": ownerName,
			"eventName": eventName,
			"eventUrl":  fmt.Sprintf("%s/e/%s", baseUrl, eventId),
		}
	} else {
		initialTmpl = "ReminderInitial"
		secondKey = "ReminderSecond"
		finalKey = "ReminderFinal"
		data = bson.M{
			"ownerName":   ownerName,
			"eventName":   eventName,
			"eventUrl":    fmt.Sprintf("%s/e/%s", baseUrl, eventId),
			"finishedUrl": fmt.Sprintf("%s/e/%s/responded?email=%s", baseUrl, eventId, email),
		}
	}

	tmpl := templateRegistry[initialTmpl]
	emailsvc.SendAndAddSubscriber(email, tmpl, toStringMap(data), false)

	id1 := insertScheduled(email, secondKey, data, time.Now().Add(24*time.Hour))
	id2 := insertScheduled(email, finalKey, data, time.Now().Add(72*time.Hour))

	ids := make([]string, 0, 2)
	if id1 != "" {
		ids = append(ids, id1)
	}
	if id2 != "" {
		ids = append(ids, id2)
	}
	return ids
}

// CancelEmail deletes a scheduled email by its ID (no-op if already sent or not found).
func CancelEmail(id string) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	db.ScheduledEmailsCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})
}

func insertScheduled(to, templateKey string, data bson.M, at time.Time) string {
	doc := scheduledEmail{
		Id:          primitive.NewObjectID(),
		ScheduledAt: at,
		To:          to,
		TemplateKey: templateKey,
		Data:        data,
	}
	_, err := db.ScheduledEmailsCollection.InsertOne(context.Background(), doc)
	if err != nil {
		logger.StdErr.Println("scheduler: failed to insert scheduled email:", err)
		return ""
	}
	return doc.Id.Hex()
}

func sendDue() {
	cursor, err := db.ScheduledEmailsCollection.Find(
		context.Background(),
		bson.M{"scheduledAt": bson.M{"$lte": time.Now()}},
	)
	if err != nil {
		logger.StdErr.Println("scheduler: failed to query due emails:", err)
		return
	}
	defer cursor.Close(context.Background())

	var jobs []scheduledEmail
	if err := cursor.All(context.Background(), &jobs); err != nil {
		logger.StdErr.Println("scheduler: failed to decode due emails:", err)
		return
	}

	for _, job := range jobs {
		tmpl, ok := templateRegistry[job.TemplateKey]
		if !ok {
			logger.StdErr.Println("scheduler: unknown template key:", job.TemplateKey)
		} else {
			emailsvc.Send(job.To, tmpl, toStringMap(job.Data))
		}
		db.ScheduledEmailsCollection.DeleteOne(context.Background(), bson.M{"_id": job.Id})
	}
}

// toStringMap converts bson.M (map[string]any) to map[string]any — same underlying type,
// but makes the call sites cleaner when mixing bson.M and map[string]any.
func toStringMap(m bson.M) map[string]any {
	return map[string]any(m)
}
