package mongodb

import (
	"context"
	"log"

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
