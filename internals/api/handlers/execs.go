package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"jsantdev.com/grpc_sm_api/internals/repositories/mongodb"
	"jsantdev.com/grpc_sm_api/pkg/utils"
	pb "jsantdev.com/grpc_sm_api/proto/gen"
)

func (s *Server) GetExecs(ctx context.Context, req *pb.GetExecsRequest) (*pb.Execs, error) {
	execFilter := req.GetExec()
	sortFields := req.GetSortBy()

	pageNumber := req.PageNumber
	pageSize := req.PageSize

	if pageNumber < 1 {
		pageNumber = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}
	execs, err := mongodb.GetExecs(ctx, execFilter, sortFields, pageNumber, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Execs{Execs: execs}, nil
}

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

func (s *Server) UpdateExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	execsFromReq := req.GetExecs()

	for _, exec := range execsFromReq {
		if exec.Id == "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: empty ID field is not allowed")
		}
	}
	updatedExecs, err := mongodb.UpdateExecs(ctx, execsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Execs{Execs: updatedExecs}, nil
}

func (s *Server) DeleteExecs(ctx context.Context, req *pb.ExecIds) (*pb.DeleteExecsConfirmation, error) {

	execIdsFromReq := req.GetIds()

	deletedIds, err := mongodb.DeleteExecs(ctx, execIdsFromReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteExecsConfirmation{
		Status:     "Execs deleted successfully",
		DeletedIds: deletedIds,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.ExecLoginRequest) (*pb.ExecLoginResponse, error) {

	username := req.GetUsername()
	password := req.GetPassword()

	if username == "" || password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password is required")
	}

	token, err := mongodb.Login(ctx, username, password)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &pb.ExecLoginResponse{
		Status: "Login successfully",
		Token:  token,
	}, nil
}

func (s *Server) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {

	userId := req.GetId()
	currentPassword := req.GetCurrentPassword()
	newPassword := req.GetNewPassword()

	if userId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty id is not allowed")
	}

	if currentPassword == "" || newPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "current/new password is required")
	}

	token, err := mongodb.UpdatePassword(ctx, userId, currentPassword, newPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePasswordResponse{
		PasswordUpdated: true,
		Token:           token,
	}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *pb.ExecIds) (*pb.Confimation, error) {

	userIdsFromReq := req.GetIds()

	if err := mongodb.DeactivateUser(ctx, userIdsFromReq); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Confimation{
		Confirmation: true,
	}, nil
}

func (s *Server) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {

	email := req.GetEmail()

	if err := mongodb.ForgotPassword(ctx, email); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ForgotPasswordResponse{
		Confirmation: true,
		Message:      "Forgot password email sent",
	}, nil

}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.Confimation, error) {

	token := req.GetResetCode()
	confirmPassword := req.GetConfirmPassword()
	newPassword := req.GetNewPassword()

	if confirmPassword == "" || newPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new/confirm password is required")
	}

	if confirmPassword != newPassword {
		return nil, status.Error(codes.InvalidArgument, "password should match")
	}

	if err := mongodb.ResetPassword(ctx, token, newPassword); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Confimation{
		Confirmation: true,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.EmptyRequest) (*pb.ExecLogoutResponsesponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata available")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	expiryTimestamp := ctx.Value(utils.ContextKey("expiresAt"))
	expiryTimestampStr := fmt.Sprintf("%v", expiryTimestamp)

	expiryTimeInt, err := strconv.ParseInt(expiryTimestampStr, 10, 64)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to extract exipiry time")
	}

	expiryTime := time.Unix(expiryTimeInt, 0)

	utils.JwtStore.AddToken(token, expiryTime)

	return &pb.ExecLogoutResponsesponse{
		LoggedOut: true,
	}, nil
}
