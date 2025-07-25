package mongodb

import (
	"context"
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
