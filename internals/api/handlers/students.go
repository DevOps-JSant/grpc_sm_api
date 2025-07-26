package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jsantdev.com/grpc_sm_api/internals/repositories/mongodb"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsRequest) (*pb.Students, error) {
	studentFilter := req.GetStudent()
	sortFields := req.GetSortBy()

	pageNumber := req.PageNumber
	pageSize := req.PageSize

	if pageNumber < 1 {
		pageNumber = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}
	students, err := mongodb.GetStudents(ctx, studentFilter, sortFields, pageNumber, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Students{Students: students}, nil
}

func (s *Server) AddStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {

	studentsFromReq := req.GetStudents()

	for _, student := range studentsFromReq {
		if student.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: non-empty ID fields are not allowed")
		}
	}

	addedStudents, err := mongodb.AddStudents(ctx, studentsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Students{Students: addedStudents}, nil

}

func (s *Server) UpdateStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {
	studentsFromReq := req.GetStudents()

	for _, student := range studentsFromReq {
		if student.Id == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: empty ID field is not allowed")
		}
	}
	updatedStudents, err := mongodb.UpdateStudents(ctx, studentsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Students{Students: updatedStudents}, nil
}

func (s *Server) DeleteStudents(ctx context.Context, req *pb.StudentIds) (*pb.DeleteStudentsConfirmation, error) {

	studentIdsFromReq := req.GetIds()

	deletedIds, err := mongodb.DeleteStudents(ctx, studentIdsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteStudentsConfirmation{
		Status:     "Students deleted successfully",
		DeletedIds: deletedIds,
	}, nil
}
