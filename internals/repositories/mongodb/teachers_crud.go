package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func GetTeachers(ctx context.Context, teacherFilterFromReq *pb.Teacher, sortFieldsFromReq []*pb.SortField) ([]*pb.Teacher, error) {
	// Filtering, getting the filters from the request
	filters, err := utils.BuildFilter(teacherFilterFromReq, &models.Teacher{})
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

	col := client.Database("school").Collection("teachers")
	var cursor *mongo.Cursor
	if len(sortOptions) > 0 {
		cursor, err = col.Find(ctx, filters, options.Find().SetSort(sortOptions))
	} else {
		cursor, err = col.Find(ctx, filters)
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to retrieve data")
	}
	defer cursor.Close(ctx)

	teachers, err := utils.DecodeEntities(ctx, cursor, func() *pb.Teacher { return &pb.Teacher{} }, func() *models.Teacher { return &models.Teacher{} })
	if err != nil {
		return nil, err
	}

	return teachers, nil
}

func AddTeachers(ctx context.Context, teachersFromReq []*pb.Teacher) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	newTeachers := make([]models.AddTeacherRequest, len(teachersFromReq))

	for i, pbTeacher := range teachersFromReq {
		modelTeacher := utils.MapPBToModel(models.AddTeacherRequest{}, pbTeacher)
		newTeachers[i] = modelTeacher
	}

	var addedTeachers []*pb.Teacher
	for _, teacher := range newTeachers {
		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to add value to database")
		}

		// pbTeacher := mapModelTeacherToPB(teacher)
		pbTeacher := utils.MapModelToPB(pb.Teacher{}, teacher)

		objectId, ok := result.InsertedID.(bson.ObjectID)
		if ok {
			pbTeacher.Id = objectId.Hex()
		}

		addedTeachers = append(addedTeachers, pbTeacher)
	}

	return addedTeachers, nil
}

func UpdateTeachers(ctx context.Context, teachersFromReq []*pb.Teacher) ([]*pb.Teacher, error) {

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

	var updatedTeachers []*pb.Teacher
	for _, teacher := range teachersFromReq {

		// Map pb.Teacher to models.UpdateTeacherRequest
		modelTeacher := utils.MapPBToModel(models.UpdateTeacherRequest{}, teacher)

		// Extract bson.ObjectId from model
		objectId, err := bson.ObjectIDFromHex(modelTeacher.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		// Convert models.UpdateTeacherRequest to BSON document
		modelDoc, err := bson.Marshal(modelTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to parse model to BSON document")
		}

		var updatedDoc bson.M
		err = bson.Unmarshal(modelDoc, &updatedDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to parse BSON document")
		}

		// Remove the _id field from the updated document
		delete(updatedDoc, "_id")

		_, err = client.Database("school").Collection("teachers").UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Unable to update teacher with id %s", modelTeacher.Id))
		}

		updatedTeacher := utils.MapModelToPB(pb.Teacher{}, modelTeacher)
		updatedTeachers = append(updatedTeachers, updatedTeacher)

	}

	return updatedTeachers, nil

}
