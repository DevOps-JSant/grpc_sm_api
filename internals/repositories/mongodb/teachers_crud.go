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
