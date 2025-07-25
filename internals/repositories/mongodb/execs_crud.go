package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func AddExecs(ctx context.Context, execsFromReq []*pb.Exec) ([]*pb.Exec, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	var addedExecs []*pb.Exec
	for _, exec := range execsFromReq {
		addExecRequest := utils.MapPBToModel(models.AddExecRequest{}, exec)

		if addExecRequest.Password == "" {
			return nil, utils.ErrorHandler(errors.New("password is required"), "Unable to add data")
		}

		// Hash password
		encodedHash, err := utils.HashPassword(exec.Password)
		if err != nil {
			return nil, err
		}
		addExecRequest.Password = encodedHash

		// Get current time
		currentTime := time.Now().Format(time.RFC3339)
		addExecRequest.UserCreatedAt = currentTime
		addExecRequest.InactiveStatus = false

		result, err := client.Database("school").Collection("execs").InsertOne(ctx, addExecRequest)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to add value to database")
		}

		pbExec := utils.MapModelToPB(pb.Exec{}, addExecRequest)
		objectId, ok := result.InsertedID.(bson.ObjectID)
		if ok {
			pbExec.Id = objectId.Hex()
		}

		addedExecs = append(addedExecs, pbExec)
	}

	return addedExecs, nil
}
