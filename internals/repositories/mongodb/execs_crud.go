package mongodb

import (
	"context"
	"errors"
	"fmt"
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

func UpdateExecs(ctx context.Context, execsFromReq []*pb.Exec) ([]*pb.Exec, error) {

	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	var updatedExecs []*pb.Exec
	for _, exec := range execsFromReq {

		// Map pb.Exec to models.UpdateExecRequest
		updateExecRequest := utils.MapPBToModel(models.UpdateExecRequest{}, exec)

		// Extract bson.ObjectId from model
		objectId, err := bson.ObjectIDFromHex(updateExecRequest.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		// Convert models.UpdateExecRequest to BSON document
		modelDoc, err := bson.Marshal(updateExecRequest)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to parse model to BSON document")
		}

		var updatedDoc bson.M
		err = bson.Unmarshal(modelDoc, &updatedDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to parse BSON document")
		}

		// Remove the _id field from the updated document
		if docPassword, ok := updatedDoc["password"]; ok {
			password := docPassword.(string)
			if password == "" {
				return nil, utils.ErrorHandler(err, fmt.Sprintf("Unable to update exec with id %s as the set password is empty", updateExecRequest.Id))
			}

			// Hash password
			encodedHash, err := utils.HashPassword(password)
			if err != nil {
				return nil, err
			}
			updatedDoc["password"] = encodedHash
			updateExecRequest.Password = encodedHash
		}

		delete(updatedDoc, "_id")

		_, err = client.Database("school").Collection("execs").UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Unable to update exec with id %s", updateExecRequest.Id))
		}

		updatedExec := utils.MapModelToPB(pb.Exec{}, updateExecRequest)
		updatedExecs = append(updatedExecs, updatedExec)

	}

	return updatedExecs, nil

}

func DeleteExecs(ctx context.Context, execIdsFromReq []*pb.ExecId) ([]string, error) {

	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	objectIds := make([]bson.ObjectID, len(execIdsFromReq))
	for i, execId := range execIdsFromReq {
		objectId, err := bson.ObjectIDFromHex(execId.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	result, err := client.Database("school").Collection("execs").DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to delete exec")
	}

	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(errors.New("no exec deleted"), "no exec deleted")
	}

	deletedIds := make([]string, result.DeletedCount)
	for i, id := range objectIds {
		deletedIds[i] = id.Hex()
	}

	return deletedIds, nil
}
