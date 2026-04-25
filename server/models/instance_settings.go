package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type InstanceSettings struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AllowRegistration bool               `json:"allowRegistration" bson:"allowRegistration"`
}
