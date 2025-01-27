package main

import (
	"context"
	"log"
	"net"

	_ "github.com/brianvoe/gofakeit"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/katyafirstova/auth_service/internal/api/user"
	userRepository "github.com/katyafirstova/auth_service/internal/repository/user"
	userService "github.com/katyafirstova/auth_service/internal/service/user"
	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

const (
	dbDSN   = "host=localhost port=54322 dbname=auth_db user=auth_user password=auth_password sslmode=disable"
	address = "127.0.0.1:50001"
)

func main() {
	var (
		pool *pgxpool.Pool
		err  error
		ctx  = context.Background()
	)

	pool, err = pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	defer pool.Close()

	userRepo := userRepository.NewRepository(pool)
	userServ := userService.NewService(userRepo)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	user_v1.RegisterUserV1Server(grpcServer, user.NewImplementation(userServ))

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
