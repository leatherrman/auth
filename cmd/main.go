package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

const (
	dbDSN = "host=localhost port=54322 dbname=auth_db user=auth_user password=auth_password sslmode=disable"

	authTable                      = "auth"
	authTableColumnId              = "id"
	authTableColumnName            = "name"
	authTableColumnEmail           = "email"
	authTableColumnPassword        = "password"
	authTableColumnPasswordConfirm = "password_confirm"
	authTableColumnRole            = "role"
	authTableColumnCreatedAt       = "created_at"
	authTableColumnUpdatedAt       = "updated_at"
)

func Create(ctx context.Context, pool *pgxpool.Pool, name string, email string, password string, role user_v1.Role) {
	builderInsert := sq.Insert(authTable).
		PlaceholderFormat(sq.Dollar).
		Columns(authTableColumnName, authTableColumnEmail, authTableColumnPassword, authTableColumnPasswordConfirm, authTableColumnRole).
		Values(name, email, password, password, role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	//if err != nil {
	//	return error(err)
	//}

	var authId int64
	err = pool.QueryRow(ctx, query, args...).Scan(&authId)
	if err != nil {
		log.Fatalf("Failed to insert to database: %s", err.Error())
	}

	log.Printf("Inserted query with id: %d", authId)
}

func Get(ctx context.Context, pool *pgxpool.Pool) {
	builderSelect := sq.Select(authTableColumnId, authTableColumnName, authTableColumnEmail, authTableColumnRole,
		authTableColumnCreatedAt, authTableColumnUpdatedAt).
		From(authTable).
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Fatalf("Failed to create query: %s", err.Error())
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Failed to select from database: %s", err.Error())
	}

	var authId int64
	var name string
	var email string
	var role user_v1.Role
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&authId, &name, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("Failed to scan from database: %s", err.Error())
		}
		log.Printf("Selected query: id: #{authId}, name: #{name}, email: #{email}, role: #{role}, " +
			"created_at: #{createdAt}, updated_at: #{updatedAt}")
	}
	defer rows.Close()
}

func Update(ctx context.Context, pool *pgxpool.Pool, id int64) {

	builderUpdate := sq.Update(authTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{authTableColumnId: id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("Failed to build query: %s", err.Error())
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("Failed to execute query: %s", err.Error())
	}

	log.Printf("Updated #{res.RowsAffected()} rows")

	log.Printf("Updated query: id: #{id}, name: #{name}, email: #{email}, role: #{role}, createdAt: " +
		"#{createdAt}, updatedAt: #{updatedAt}")
}

func DeleteById(ctx context.Context, pool *pgxpool.Pool, id int64) {
	builderDelete := sq.Delete(authTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{authTableColumnId: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("Failed to build query: %s", err.Error())
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("Failed to execute query: %s", err.Error())
	}

	log.Printf("Delete query with id: #{authId}")
	log.Printf("Updated #{res.RowsAffected()} rows")
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	defer pool.Close()

	Create(ctx, pool, gofakeit.Word(), gofakeit.Email(), gofakeit.Word(), user_v1.Role_USER)
	Get(context.Background(), pool)
	Update(context.Background(), pool, 1)
	DeleteById(context.Background(), pool, 2)
}
