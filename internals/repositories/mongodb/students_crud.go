package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func GetStudents(ctx context.Context, studentFilterFromReq *pb.Student, sortFieldsFromReq []*pb.SortField, pageNumber, pageSize uint32) ([]*pb.Student, error) {
	// Filtering, getting the filters from the request
	filters, err := utils.BuildFilter(studentFilterFromReq, &models.Student{})
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

	col := client.Database("school").Collection("students")

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

	students, err := utils.DecodeEntities(ctx, cursor, func() *pb.Student { return &pb.Student{} }, func() *models.Student { return &models.Student{} })
	if err != nil {
		return nil, err
	}

	return students, nil
}

func AddStudents(ctx context.Context, studentsFromReq []*pb.Student) ([]*pb.Student, error) {
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	var addedStudents []*pb.Student
	for _, student := range studentsFromReq {
		addStudentRequest := utils.MapPBToModel(models.AddStudentRequest{}, student)

		result, err := client.Database("school").Collection("students").InsertOne(ctx, addStudentRequest)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to add value to database")
		}

		pbStudent := utils.MapModelToPB(pb.Student{}, addStudentRequest)
		objectId, ok := result.InsertedID.(bson.ObjectID)
		if ok {
			pbStudent.Id = objectId.Hex()
		}

		addedStudents = append(addedStudents, pbStudent)
	}

	return addedStudents, nil
}

func UpdateStudents(ctx context.Context, studentsFromReq []*pb.Student) ([]*pb.Student, error) {

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

	var updatedStudents []*pb.Student
	for _, student := range studentsFromReq {

		// Map pb.Teacher to models.UpdateStudentRequest
		updateStudentRequest := utils.MapPBToModel(models.UpdateStudentRequest{}, student)

		// Extract bson.ObjectId from model
		objectId, err := bson.ObjectIDFromHex(updateStudentRequest.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		// Convert models.UpdateStudentRequest to BSON document
		modelDoc, err := bson.Marshal(updateStudentRequest)
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

		_, err = client.Database("school").Collection("students").UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Unable to update student with id %s", updateStudentRequest.Id))
		}

		updatedStudent := utils.MapModelToPB(pb.Student{}, updateStudentRequest)
		updatedStudents = append(updatedStudents, updatedStudent)

	}

	return updatedStudents, nil

}

func DeleteStudents(ctx context.Context, studentIdsFromReq []*pb.StudentId) ([]string, error) {

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

	objectIds := make([]bson.ObjectID, len(studentIdsFromReq))
	for i, studentId := range studentIdsFromReq {
		objectId, err := bson.ObjectIDFromHex(studentId.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	result, err := client.Database("school").Collection("students").DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to delete student")
	}

	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(errors.New("no student deleted"), "no student deleted")
	}

	deletedIds := make([]string, result.DeletedCount)
	for i, id := range objectIds {
		deletedIds[i] = id.Hex()
	}

	return deletedIds, nil
}
