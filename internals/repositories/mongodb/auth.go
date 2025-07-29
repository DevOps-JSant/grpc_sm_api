package mongodb

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	"jsantdev.com/grpc_sm_api/pkg/utils/templates"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
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
		return "", utils.ErrorHandler(errors.New("user is inactive"), "user is inactive")
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

func DeactivateUser(ctx context.Context, userIdsFromReq []*pb.ExecId) error {
	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	objectIds := make([]bson.ObjectID, len(userIdsFromReq))
	for i, exec := range userIdsFromReq {

		// Extract bson.ObjectId from model
		objectId, err := bson.ObjectIDFromHex(exec.Id)
		if err != nil {
			return utils.ErrorHandler(err, "Invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	update := bson.M{
		"$set": bson.M{
			"inactive_status": true,
		},
	}
	_, err = client.Database("school").Collection("execs").UpdateMany(ctx, filter, update)
	if err != nil {
		return utils.ErrorHandler(err, "unable to deactive user")
	}

	return nil
}

func ForgotPassword(ctx context.Context, emailFromReq string) error {
	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	var user models.Exec
	err = client.Database("school").Collection("execs").FindOne(ctx, bson.M{"email": emailFromReq}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorHandler(err, "user not found")
		}
		return utils.ErrorHandler(err, "internal error")
	}

	duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		return utils.ErrorHandler(err, "unable to send reset password")
	}

	mins := time.Duration(duration)

	expiry := time.Now().Add(mins * time.Minute).Format(time.RFC3339)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return utils.ErrorHandler(err, "unable to send reset password")
	}

	token := hex.EncodeToString(tokenBytes)
	hashedToken := sha256.Sum256(tokenBytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	update := bson.M{
		"$set": bson.M{
			"password_reset_token":   hashedTokenString,
			"password_token_expires": expiry,
		},
	}
	filter := bson.M{
		"_id": user.Id,
	}

	_, err = client.Database("school").Collection("execs").UpdateOne(ctx, filter, update)
	if err != nil {
		return utils.ErrorHandler(err, "unable to send reset password")
	}

	// Send email
	resetURL := fmt.Sprintf("https://localhost:3000/resetpassword/reset/%s", token)

	// Parse the template file
	tmpl, err := template.ParseFS(templates.TemplateFS, "reset_password.html")
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send reset password")
	}

	// place holder data
	placeHolderData := struct {
		Username  string
		ResetLink string
		Expiry    int
	}{
		Username:  user.Username,
		ResetLink: resetURL,
		Expiry:    int(mins),
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, placeHolderData)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send reset password")
	}

	message := buf.String()

	emailSender := utils.Email{
		From:     "support@jsantdev.com",
		To:       user.Email,
		Subject:  "Your password reset link",
		BodyType: "text/html",
		Body:     message,
	}

	if err := emailSender.Send(); err != nil {
		cleanup := bson.M{
			"$set": bson.M{
				"password_reset_token":   nil,
				"password_token_expires": nil,
			},
		}

		_, _ = client.Database("school").Collection("execs").UpdateOne(ctx, filter, cleanup)
		return err
	}

	return nil
}

func ResetPassword(ctx context.Context, token, newPassword string) error {

	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	// decode the token
	bytes, err := hex.DecodeString(token)
	if err != nil {
		return utils.ErrorHandler(err, "Internal error")
	}

	hashedToken := sha256.Sum256(bytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	// Get user by password_reset_token
	filter := bson.M{
		"password_reset_token": hashedTokenString,
		"password_token_expires": bson.M{
			"$gt": time.Now().Format(time.RFC3339),
		},
	}

	var user models.Exec
	err = client.Database("school").Collection("execs").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return utils.ErrorHandler(err, "Invalid token or token expires")
	}

	// hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the password
	update := bson.M{
		"$set": bson.M{
			"password":               hashedPassword,
			"password_changed_at":    time.Now().Format(time.RFC3339),
			"password_reset_token":   nil,
			"password_token_expires": nil,
		},
	}
	_, err = client.Database("school").Collection("execs").UpdateOne(ctx, bson.M{"_id": user.Id}, update)
	if err != nil {
		return utils.ErrorHandler(err, "Unable to reset password")
	}

	return nil
}
