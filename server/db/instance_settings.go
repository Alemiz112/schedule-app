package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/logger"
	"schej.it/server/models"
)

func GetInstanceSettings() models.InstanceSettings {
	var settings models.InstanceSettings
	err := InstanceSettingsCollection.FindOne(context.Background(), bson.M{}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		return models.InstanceSettings{AllowRegistration: true}
	}
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return settings
}

func UpsertInstanceSettings(settings models.InstanceSettings) error {
	opts := options.Replace().SetUpsert(true)
	_, err := InstanceSettingsCollection.ReplaceOne(
		context.Background(),
		bson.M{},
		settings,
		opts,
	)
	return err
}
