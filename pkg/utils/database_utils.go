package utils

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
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
