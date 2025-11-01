package main

import (
	"context"
	"fmt"

	pb "github.com/Riku-KANO/kube-ec/services/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	repo *UserRepository
}

func NewUserServer(repo *UserRepository) *UserServer {
	return &UserServer{
		repo: repo,
	}
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userRow, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return userRowToProto(userRow), nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.repo.Update(ctx, req.Id, req.Name, req.PhoneNumber); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	userRow, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated user")
	}

	return userRowToProto(userRow), nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.repo.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete user: %v", err))
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}
