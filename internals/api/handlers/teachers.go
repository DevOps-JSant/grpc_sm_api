package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jsantdev.com/grpc_sm_api/internals/repositories/mongodb"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.Teachers, error) {

	teacherFilter := req.GetTeacher()
	sortFields := req.GetSortBy()

	teachers, err := mongodb.GetTeachers(ctx, teacherFilter, sortFields)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: teachers}, nil
}

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	teachersFromReq := req.GetTeachers()

	for _, teacher := range teachersFromReq {
		if teacher.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: non-empty ID fields are not allowed")
		}
	}

	addedTeachers, err := mongodb.AddTeachers(ctx, teachersFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Teachers{Teachers: addedTeachers}, nil

}

func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	teachersFromReq := req.GetTeachers()

	for _, teacher := range teachersFromReq {
		if teacher.Id == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: empty ID field is not allowed")
		}
	}
	updatedTeachers, err := mongodb.UpdateTeachers(ctx, teachersFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Teachers{Teachers: updatedTeachers}, nil
}
