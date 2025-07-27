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

	var addedTeachers []*pb.Teacher
	for _, teacher := range teachersFromReq {
		addTeacherRequest := utils.MapPBToModel(models.AddTeacherRequest{}, teacher)

		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, addTeacherRequest)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to add value to database")
		}

		pbTeacher := utils.MapModelToPB(pb.Teacher{}, addTeacherRequest)
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
		updateTeacherRequest := utils.MapPBToModel(models.UpdateTeacherRequest{}, teacher)

		// Extract bson.ObjectId from model
		objectId, err := bson.ObjectIDFromHex(updateTeacherRequest.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}

		// Convert models.UpdateTeacherRequest to BSON document
		modelDoc, err := bson.Marshal(updateTeacherRequest)
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
			return nil, utils.ErrorHandler(err, fmt.Sprintf("Unable to update teacher with id %s", updateTeacherRequest.Id))
		}

		updatedTeacher := utils.MapModelToPB(pb.Teacher{}, updateTeacherRequest)
		updatedTeachers = append(updatedTeachers, updatedTeacher)

	}

	return updatedTeachers, nil

}

func DeleteTeachers(ctx context.Context, teacherIdsFromReq []*pb.TeacherId) ([]string, error) {

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

	objectIds := make([]bson.ObjectID, len(teacherIdsFromReq))
	for i, teacherId := range teacherIdsFromReq {
		objectId, err := bson.ObjectIDFromHex(teacherId.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	result, err := client.Database("school").Collection("teachers").DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Unable to delete teacher")
	}

	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(errors.New("no teachers deleted"), "no teachers deleted")
	}

	deletedIds := make([]string, result.DeletedCount)
	for i, id := range objectIds {
		deletedIds[i] = id.Hex()
	}

	return deletedIds, nil
}

func GetStudentCountByClassTeacher(ctx context.Context, teacherIdFromReq string) (int, error) {
	// Connect to mongo db
	client, err := CreateMongoClient(ctx)
	if err != nil {
		return 0, err
	}
	// Close connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Unable to disconnect to mongodb:", err)
		}
	}()

	objectId, err := bson.ObjectIDFromHex(teacherIdFromReq)
	if err != nil {
		return 0, utils.ErrorHandler(err, "Invalid id")
	}

	var teacher models.Teacher
	err = client.Database("school").Collection("teachers").FindOne(ctx, bson.M{"_id": objectId}).Decode(&teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, utils.ErrorHandler(err, "No teacher found")
		}
		return 0, utils.ErrorHandler(err, "Internal error")
	}

	studentCount, err := client.Database("school").Collection("students").CountDocuments(ctx, bson.M{"class": teacher.Class})
	if err != nil {
		return 0, utils.ErrorHandler(err, "Unable to count student by class teacher")
	}

	return int(studentCount), nil
}

func GetStudentsByClassTeacher(ctx context.Context, teacherIdFromReq string) ([]*pb.Student, error) {
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

	objectId, err := bson.ObjectIDFromHex(teacherIdFromReq)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Invalid id")
	}

	var teacher models.Teacher
	err = client.Database("school").Collection("teachers").FindOne(ctx, bson.M{"_id": objectId}).Decode(&teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrorHandler(err, "No teacher found")
		}
		return nil, utils.ErrorHandler(err, "Internal error")
	}

	col := client.Database("school").Collection("students")
	var cursor *mongo.Cursor
	cursor, err = col.Find(ctx, bson.M{"class": teacher.Class})
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
