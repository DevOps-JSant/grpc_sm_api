package utils

import (
	"context"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func BuildFilter(pbModel interface{}, model interface{}) (bson.M, error) {
	filter := bson.M{}

	if pbModel == nil || reflect.ValueOf(pbModel).IsNil() {
		return filter, nil
	}

	// Mapping pbObject fields to model
	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	reqVal := reflect.ValueOf(pbModel).Elem()
	reqType := reqVal.Type()

	for i := range reqVal.NumField() {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				if fieldName == "Id" {
					objectId, err := bson.ObjectIDFromHex(reqVal.FieldByName(fieldName).Interface().(string))
					if err != nil {
						return nil, ErrorHandler(err, "Invalid id")
					}
					modelField.Set(reflect.ValueOf(objectId))
				} else {
					modelField.Set(fieldVal)
				}
			}
		}
	}

	// Iterate over the modelTeacher to build filter using bson.M
	for i := range modelVal.NumField() {
		fieldVal := modelVal.Field(i)
		// fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objectId, err := bson.ObjectIDFromHex(fieldVal.Interface().(bson.ObjectID).Hex())
				if err != nil {
					return nil, ErrorHandler(err, "Invalid id")
				}
				filter[bsonTag] = objectId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}
	return filter, nil
}

func BuildSortOptions(sortFields []*pb.SortField) bson.D {
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

func DecodeEntities[M interface{}, T interface{}](ctx context.Context, cursor *mongo.Cursor, newEntity func() *T, newModel func() *M) ([]*T, error) {
	var entities []*T
	for cursor.Next(ctx) {
		model := newModel()
		err := cursor.Decode(&model)
		if err != nil {
			return nil, ErrorHandler(err, "Unable to decode data")
		}

		entity := newEntity()
		modelVal := reflect.ValueOf(model).Elem()
		pbVal := reflect.ValueOf(entity).Elem()

		for i := range modelVal.NumField() {
			modelField := modelVal.Field(i)
			modelFieldName := modelVal.Type().Field(i).Name

			pbField := pbVal.FieldByName(modelFieldName)
			if pbField.IsValid() && pbField.CanSet() {
				if modelFieldName == "Id" {
					objectId := modelVal.FieldByName(modelFieldName).Interface().(bson.ObjectID).Hex()
					// if err != nil {
					// 	return nil, utils.ErrorHandler(err, "Invalid id")
					// }
					pbField.Set(reflect.ValueOf(objectId))
				} else {
					pbField.Set(modelField)
				}
			}
		}
		entities = append(entities, entity)
	}
	err := cursor.Err()
	if err != nil {
		return nil, ErrorHandler(err, "Error in cursor")
	}
	return entities, nil
}

// Map model to pb entities
func MapModelToPB[T1 any, T2 any](target T1, source T2) *T1 {
	pbResult := &target

	modelVal := reflect.ValueOf(source)
	pbVal := reflect.ValueOf(pbResult).Elem()

	for i := range modelVal.NumField() {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbResult

}

// Map pb entities to model
func MapPBToModel[T1 any, T2 any](target T1, source T2) T1 {
	modelResult := target

	pbVal := reflect.ValueOf(source).Elem()
	modelVal := reflect.ValueOf(&modelResult).Elem()

	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return modelResult
}
