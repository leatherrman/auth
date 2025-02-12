package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/katyafirstova/auth_service/internal/model"
	"github.com/katyafirstova/auth_service/internal/repository/user/converter"
	modelRepo "github.com/katyafirstova/auth_service/internal/repository/user/model"
)

const (
	userTable                = "users"
	userTableColumnUUID      = "uuid"
	userTableColumnName      = "name"
	userTableColumnEmail     = "email"
	userTableColumnPassword  = "password"
	userTableColumnRole      = "role"
	userTableColumnCreatedAt = "created_at"
	userTableColumnUpdatedAt = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, req model.CreateUser) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(userTableColumnUUID, userTableColumnName, userTableColumnEmail, userTableColumnPassword, userTableColumnRole).
		Values(uuid.NewString(), req.Name, req.Email, hashedPassword, req.Role).
		Suffix(fmt.Sprintf("RETURNING %s", userTableColumnUUID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return "", err
	}

	var newUUID string
	err = r.db.QueryRow(ctx, query, args...).Scan(&newUUID)
	if err != nil {
		return "", err
	}

	return newUUID, nil
}

func (r *repo) Get(ctx context.Context, uuid string) (model.User, error) {
	builderSelect := sq.Select(userTableColumnUUID, userTableColumnName, userTableColumnEmail, userTableColumnRole,
		userTableColumnCreatedAt, userTableColumnUpdatedAt).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		OrderBy(fmt.Sprintf("%s ASC", userTableColumnUUID)).
		Limit(1).
		Where(sq.Eq{userTableColumnUUID: uuid})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user modelRepo.User

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.UUID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return converter.ToUserFromRepo(user), nil
}

func (r *repo) Update(ctx context.Context, uuid string, req model.UpdateUser) error {
	builderUpdate := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set(userTableColumnUpdatedAt, time.Now()).
		Where(sq.Eq{userTableColumnUUID: uuid})

	if req.Name != nil {
		builderUpdate = builderUpdate.Set(userTableColumnName, req.Name)
	}

	if req.Email != nil {
		builderUpdate = builderUpdate.Set(userTableColumnEmail, req.Email)
	}

	if req.Role != 0 {
		builderUpdate = builderUpdate.Set(userTableColumnRole, req.Role)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, uuid string) error {
	builderDelete := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userTableColumnUUID: uuid})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
