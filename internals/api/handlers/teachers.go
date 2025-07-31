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

	err := req.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request %v", err)
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

func (s *Server) DeleteTeachers(ctx context.Context, req *pb.TeacherIds) (*pb.DeleteTeachersConfirmation, error) {

	teacherIdsFromReq := req.GetIds()

	deletedIds, err := mongodb.DeleteTeachers(ctx, teacherIdsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteTeachersConfirmation{
		Status:     "Teachers deleted successfully",
		DeletedIds: deletedIds,
	}, nil
}

func (s *Server) GetStudentCountByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.StudentCount, error) {

	teacherIdFromReq := req.GetId()

	if teacherIdFromReq == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	studentCount, err := mongodb.GetStudentCountByClassTeacher(ctx, teacherIdFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.StudentCount{
		StudentCount: int32(studentCount),
		Status:       true,
	}, nil

}

func (s *Server) GetStudentsByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.Students, error) {

	teacherIdFromReq := req.GetId()

	if teacherIdFromReq == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid id")
	}

	students, err := mongodb.GetStudentsByClassTeacher(ctx, teacherIdFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Students{
		Students: students,
	}, nil
}
