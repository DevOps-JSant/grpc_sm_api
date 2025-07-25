package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func GetExecs(ctx context.Context, execFilterFromReq *pb.Exec, sortFieldsFromReq []*pb.SortField, pageNumber, pageSize uint32) ([]*pb.Exec, error) {
	// Filtering, getting the filters from the request
	filters, err := utils.BuildFilter(execFilterFromReq, &models.Exec{})
	if err != nil {
		return nil, err
	}

	// Sorting. getting the sort options from the request
	sortOptions := utils.BuildSortOptions(sortFieldsFromReq)

	// Access the database to fetch data
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	col := client.Database("school").Collection("execs")

	findOptions := options.Find()
	findOptions.SetSkip(int64((pageNumber - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	var cursor *mongo.Cursor
	if len(sortOptions) > 0 {
		findOptions.SetSort(sortOptions)
	}

	cursor, err = col.Find(ctx, filters, findOptions)

	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to retrieve data")
	}
	defer cursor.Close(ctx)

	execs, err := utils.DecodeEntities(ctx, cursor, func() *pb.Exec { return &pb.Exec{} }, func() *models.Exec { return &models.Exec{} })
	if err != nil {
		return nil, err
	}

	return execs, nil
}

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
