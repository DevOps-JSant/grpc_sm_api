package mongodb

import (
	"context"
	"log"
	"reflect"
	"strings"

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
		modelTeacher := mapPBTeacherToModel(pbTeacher)
		newTeachers[i] = modelTeacher
	}

	var addedTeachers []*pb.Teacher
	for _, teacher := range newTeachers {
		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to add value to database")
		}

		pbTeacher := mapModelTeacherToPB(teacher)

		objectId, ok := result.InsertedID.(bson.ObjectID)
		if ok {
			pbTeacher.Id = objectId.Hex()
		}

		addedTeachers = append(addedTeachers, pbTeacher)
	}

	return addedTeachers, nil
}

func mapModelTeacherToPB(teacher models.AddTeacherRequest) *pb.Teacher {
	pbTeacher := &pb.Teacher{}
	modelVal := reflect.ValueOf(teacher)
	pbVal := reflect.ValueOf(pbTeacher).Elem()

	for i := range modelVal.NumField() {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbTeacher
}

func mapPBTeacherToModel(pbTeacher *pb.Teacher) models.AddTeacherRequest {
	modelTeacher := models.AddTeacherRequest{}

	pbVal := reflect.ValueOf(pbTeacher).Elem()
	modelVal := reflect.ValueOf(&modelTeacher).Elem()

	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return modelTeacher
}

func GetTeachers(ctx context.Context, teacherFilterFromReq *pb.Teacher, sortFieldsFromReq []*pb.SortField) ([]*pb.Teacher, error) {
	// Filtering, getting the filters from the request
	filters, err := buildFilterForTeachers(teacherFilterFromReq)
	if err != nil {
		return nil, err
	}

	// Sorting. getting the sort options from the request
	sortOptions := buildSortOptions(sortFieldsFromReq)

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

	var teachers []*pb.Teacher
	for cursor.Next(ctx) {
		var teacher models.Teacher
		err := cursor.Decode(&teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Unable to decode data")
		}
		teachers = append(teachers, &pb.Teacher{
			Id:        teacher.Id.Hex(),
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Email:     teacher.Email,
			Class:     teacher.Class,
			Subject:   teacher.Subject,
		})
	}

	return teachers, nil
}

func buildFilterForTeachers(teacherFilter *pb.Teacher) (bson.M, error) {
	filter := bson.M{}

	if teacherFilter == nil {
		return filter, nil
	}

	// Mapping pb.Teacher fields to models.Teacher
	var modelTeacher models.Teacher
	modelVal := reflect.ValueOf(&modelTeacher).Elem()
	modelType := modelVal.Type()

	// reqTeacher := req.GetTeacher()
	reqVal := reflect.ValueOf(teacherFilter).Elem()
	reqType := reqVal.Type()

	for i := range reqVal.NumField() {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				if fieldName == "Id" {
					objectId, err := bson.ObjectIDFromHex(teacherFilter.Id)
					if err != nil {
						return nil, utils.ErrorHandler(err, "Invalid id")
					}
					modelTeacher.Id = objectId
				} else {
					modelField.Set(fieldVal)
				}
			}
		}
	}

	// Iterate over the modelTeacher to build filter using bson.M
	for i := range modelVal.NumField() {
		fieldVal := modelVal.Field(i)

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objectId, err := bson.ObjectIDFromHex(teacherFilter.Id)
				if err != nil {
					return nil, utils.ErrorHandler(err, "Invalid id")
				}
				filter[bsonTag] = objectId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	return filter, nil
}

func buildSortOptions(sortFields []*pb.SortField) bson.D {
	var sortOptions bson.D

	for _, sortField := range sortFields {
		order := 1
		if sortField.GetOrder() == pb.Order_DESC {
			order = -1
		}

		sortOptions = append(sortOptions, bson.E{Key: sortField.Field, Value: order})
	}
	return sortOptions
}
