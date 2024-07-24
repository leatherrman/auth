package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

type server struct {
	user_v1.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	fmt.Printf("%#v", req)
	return &user_v1.CreateResponse{Id: 18}, nil
}

func (s *server) Update(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("%#v", req)
	return nil, nil
}

func (s *server) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	fmt.Printf("%#v", req)
	return &user_v1.GetResponse{Id: 18, Name: "katya", Email: "abcde@email.ru", Role: user_v1.Role_USER,
		CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now()}, nil
}

func (s *server) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("%#v", req)
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	user_v1.RegisterUserV1Server(grpcServer, &server{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
