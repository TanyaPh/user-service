package grpc_server

import (
	"context"
	"fmt"
	"log"
	"user-service/internal/models"
	"user-service/internal/services"

	pb "user-service/proto"

	"google.golang.org/grpc"
)

type UserserviceServer interface {
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(context.Context, *pb.GetUserRequest) (*pb.GetUserResponse, error)
	ListUsers(*pb.ListUsersRequest, grpc.ServerStreamingServer[pb.ListUsersResponse]) error
	UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
}

type Server struct {
	service *services.User
}

func NewServer(service *services.User) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := models.NewUser(req.Name, req.Address, req.Phone)
	userId, err := s.service.CreateUser(user)
	if err != nil {
		log.Println("couldn't add user. Error: %w", err)
		return nil, fmt.Errorf("couldn't add user")
	}

	return &pb.CreateUserResponse{ Id: userId }, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.service.GetUserById(req.Id)
	if err != nil {
		log.Println("couldn't get user. Error: %w", err)
		return nil, fmt.Errorf("couldn't get user")
	}

	return &pb.GetUserResponse {
		Id: user.Id,
		Name: user.Name,
		Address: user.Address,
		Phone: user.Phone,
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

func (s *Server) ListUsers(req *pb.ListUsersRequest, stream grpc.ServerStreamingServer[pb.ListUsersResponse]) error {
	users, err := s.service.GetUsers()
	if err != nil {
		log.Println("couldn't get users. Error: %w", err)
		return fmt.Errorf("couldn't get users")
	}

	for _, user := range users {
		err := stream.Send(&pb.ListUsersResponse {
			Id: user.Id,
			Name: user.Name,
			Address: user.Address,
			Phone: user.Phone,
			CreatedAt: user.CreatedAt.Unix(),
		})
		if err != nil {
			log.Println("couldn't send user. Error: %w", err)
		}
	}
	return nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := s.service.EditUser(req.Id, &req.Name, &req.Address, &req.Phone)
	if err != nil {
		log.Println("couldn't update user. Error: %w", err)
		return nil, fmt.Errorf("couldn't update user")
	}
	return &pb.UpdateUserResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.service.DeleteUser(req.Id)
	if err != nil {
		log.Println("couldn't delete user. Error: %w", err)
		return nil, fmt.Errorf("couldn't delete user")
	}
	return &pb.DeleteUserResponse{}, nil
}
