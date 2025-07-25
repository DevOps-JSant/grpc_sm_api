package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jsantdev.com/grpc_sm_api/internals/repositories/mongodb"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func (s *Server) AddExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {

	execsFromReq := req.GetExecs()

	for _, exec := range execsFromReq {
		if exec.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: non-empty ID fields are not allowed")
		}
	}

	addedExecs, err := mongodb.AddExecs(ctx, execsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Execs{Execs: addedExecs}, nil

}
