package handlers

import pb "jsantdev.com/grpc_sm_api/proto/gen"

type Server struct {
	pb.UnimplementedTeacherServiceServer
	pb.UnimplementedStudentServiceServer
	pb.UnimplementedExecServiceServer
}
