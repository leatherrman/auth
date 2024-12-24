package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Manzo48/microservices_course_boilerplate/pkg/users" // Замените на актуальный путь к сгенерированному gRPC коду
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServiceServer реализует интерфейс pb.UserServiceServer
type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
}

// CreateUser создает нового пользователя
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    fmt.Printf("CreateUser called with: %+v\n", req)
    return &pb.CreateUserResponse{
        Id: time.Now().Unix(), // Имитация ID
    }, nil
}

// GetUser получает информацию о пользователе по ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    fmt.Printf("GetUser called with: %+v\n", req.GetId())
    return &pb.GetUserResponse{
        Id:        req.GetId(),
        Name:      "Test User",
        Email:     "test@example.com",
        Role:      pb.Role_USER,
        CreatedAt: timestamppb.New(time.Now().Add(-time.Hour)), // Имитация времени
        UpdatedAt: timestamppb.New(time.Now()),
    }, nil
}

// UpdateUser обновляет информацию о пользователе
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
    fmt.Printf("UpdateUser called with: %+v\n", req)
    return &emptypb.Empty{}, nil
}

// DeleteUser удаляет пользователя
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
    fmt.Printf("DeleteUser called with: %+v\n", req)
    return &emptypb.Empty{}, nil
}

func main() {
    // Создание gRPC-сервера
    server := grpc.NewServer()
    pb.RegisterUserServiceServer(server, &UserServiceServer{})
	reflection.Register(server)
    // Настройка адреса
    address := ":50051"
    listener, err := net.Listen("tcp", address)
    if err != nil {
        log.Fatalf("Failed to listen on %s: %v", address, err)
    }

    log.Printf("gRPC server is running on %s", address)
    if err := server.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}