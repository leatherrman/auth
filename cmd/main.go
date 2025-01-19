package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

const (
	dbDSN   = "host=localhost port=54322 dbname=auth_db user=auth_user password=auth_password sslmode=disable"
	address = "127.0.0.1:50001"

	authTable                = "auth"
	authTableColumnID        = "id"
	authTableColumnName      = "name"
	authTableColumnEmail     = "email"
	authTableColumnPassword  = "password"
	authTableColumnRole      = "role"
	authTableColumnCreatedAt = "created_at"
	authTableColumnUpdatedAt = "updated_at"
)

var pool *pgxpool.Pool

type server struct {
	user_v1.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	if req.Password != req.PasswordConfirm {
		return nil, errors.New("passwords are not equal")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	builderInsert := sq.Insert(authTable).
		PlaceholderFormat(sq.Dollar).
		Columns(authTableColumnName, authTableColumnEmail, authTableColumnPassword, authTableColumnRole).
		Values(req.Name, req.Email, hashedPassword, req.Role).
		Suffix(fmt.Sprintf("RETURNING %s", authTableColumnID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	var id int64
	err = pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &user_v1.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	builderSelect := sq.Select(authTableColumnID, authTableColumnName, authTableColumnEmail, authTableColumnRole,
		authTableColumnCreatedAt, authTableColumnUpdatedAt).
		From(authTable).
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(1).
		Where(sq.Eq{authTableColumnID: req.Id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var id int64
	var name string
	var email string
	var role user_v1.Role
	var createdAt time.Time
	var updatedAt *time.Time

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var updatedAtTs *timestamppb.Timestamp
	if updatedAt != nil {
		updatedAtTs = timestamppb.New(*updatedAt)
	}

	return &user_v1.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtTs,
	}, nil
}

func (s *server) Update(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.Update(authTable).
		PlaceholderFormat(sq.Dollar).
		Set(authTableColumnUpdatedAt, time.Now()).
		Where(sq.Eq{authTableColumnID: req.Id})

	if req.Name != nil {
		builderUpdate = builderUpdate.Set(authTableColumnName, req.Name.Value)
	}

	if req.Email != nil {
		builderUpdate = builderUpdate.Set(authTableColumnEmail, req.Email.Value)
	}

	if req.Role != user_v1.Role_UNKNOWN {
		builderUpdate = builderUpdate.Set(authTableColumnRole, req.Role)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (s *server) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete(authTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{authTableColumnID: req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func main() {
	var err error
	ctx := context.Background()

	pool, err = pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", address)
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
