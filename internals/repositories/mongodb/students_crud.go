package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"jsantdev.com/grpc_sm_api/internals/models"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

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
