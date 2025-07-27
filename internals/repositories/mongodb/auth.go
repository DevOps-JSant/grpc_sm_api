package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
)

func Login(ctx context.Context, username, password string) (string, error) {
	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return "", err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	// Get user by username
	var user models.Exec
	filter := bson.M{"username": username}
	err = client.Database("school").Collection("execs").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", utils.ErrorHandler(err, "user not found")
		}
		return "", utils.ErrorHandler(err, "internal error")
	}

	// Check if user is inactive
	if user.InactiveStatus {
		return "", utils.ErrorHandler(err, "user is inactive")
	}

	// Verify password
	err = utils.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	uID := user.Id.Hex()
	userName := user.Username
	email := user.Email
	role := user.Role

	token, err := utils.GenerateToken(uID, userName, email, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func UpdatePassword(ctx context.Context, userId, currentPassword, newPassword string) (string, error) {
	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return "", err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	// Get user by id
	var user models.Exec
	// Extract bson.ObjectId from model
	objectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return "", utils.ErrorHandler(err, "Invalid id")
	}
	filter := bson.M{"_id": objectId}
	err = client.Database("school").Collection("execs").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", utils.ErrorHandler(err, "user not found")
		}
		return "", utils.ErrorHandler(err, "internal error")
	}

	err = utils.VerifyPassword(currentPassword, user.Password)
	if err != nil {
		return "", err
	}

	// Hash password
	encodedHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"password":            encodedHash,
			"password_changed_at": time.Now().Format(time.RFC3339),
		},
	}

	// Update password
	_, err = client.Database("school").Collection("execs").UpdateOne(ctx, filter, update)
	if err != nil {
		return "", utils.ErrorHandler(err, "unable to update password")
	}

	// Generate new token
	uID := user.Id.Hex()
	userName := user.Username
	email := user.Email
	role := user.Role

	token, err := utils.GenerateToken(uID, userName, email, role)
	if err != nil {
		return "", err
	}

	return token, nil

}
