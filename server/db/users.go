package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/logger"
	"schej.it/server/models"
)

// Returns a user based on their _id
func GetUserById(userId string) *models.User {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// userId is malformatted
		return nil
	}
	result := UsersCollection.FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// User does not exist!
		return nil
	}

	// Decode result
	var user models.User
	if err := result.Decode(&user); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &user
}


func GetAllUsers() []models.User {
	cursor, err := UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		logger.StdErr.Panicln(err)
	}
	return users
}

func GetUserCount() int64 {
	count, err := UsersCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return count
}

func UpdateUserRole(userId string, role models.UserRole) error {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	_, err = UsersCollection.UpdateByID(context.Background(), objectId, bson.M{
		"$set": bson.M{"role": role},
	})
	return err
}

func DeleteUserById(userId string) error {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	_, err = UsersCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	return err
}

func GetUserByEmail(email string) *models.User {
	result := UsersCollection.FindOne(context.Background(), bson.M{
		"email": email,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// User does not exist!
		return nil
	}

	// Decode result
	var user models.User
	if err := result.Decode(&user); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &user
}
