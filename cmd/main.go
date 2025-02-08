package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	baseURL        = "localhost:8081"
	usersPostfix   = "/users"
	userPostfix    = usersPostfix + "/{id}"
	defaultTimeout = time.Second * 5
)

// Role is ...
type Role uint8

const (
	// AdminRole is ...
	AdminRole Role = iota + 1
	// UserRole is ...
	UserRole
)

// NewUserData is ...
type NewUserData struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            uint8  `json:"role"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// UserData is ...
type UserData struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateUserData is ...
type UpdateUserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

// ResponseUserID is ...
type ResponseUserID struct {
	ID int `json:"id"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := &NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	id, errID := createUser(newUser)

	if errID != nil {
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
	}

	res := ResponseUserID{ID: id}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode new user id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func createUser(user *NewUserData) (int, error) {
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return 0, errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return 0, err
	}
	defer con.Close(ctx)

	builderInsert := sq.Insert("auth").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role", "password").
		Values(user.Name, user.Email, user.Role, user.Password).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int
	err = con.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := getUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode new user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getUser(id int) (*UserData, error) {
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return nil, errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return nil, err
	}
	defer con.Close(ctx)

	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(sq.Dollar).
		Where("id = $1", id)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var user = &UserData{}
	err = con.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updatedUser := &UpdateUserData{}
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(updatedUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	updatedUser.ID = userID
	updateUser(updatedUser)
	res := ResponseUserID{ID: userID}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode updated user id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func updateUser(user *UpdateUserData) {
	fmt.Printf("update user data: %+v\n", *user)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = deleteUser(userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(id int) error {
	pgDSN, ok := os.LookupEnv("PG_DSN")
	if !ok {
		return errors.New("PG_DSN environment variable not set")
	}

	ctx := context.Background()
	con, err := pgx.Connect(ctx, pgDSN)
	if err != nil {
		return err
	}
	defer con.Close(ctx)

	builderDelete := sq.Delete("auth").Where("id = $1", id)

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	_, err = con.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func parseID(idStr string) (int, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func main() {
	getEnv()
	r := chi.NewRouter()
	r.Post(usersPostfix, createUserHandler)
	r.Get(userPostfix, getUserHandler)
	r.Put(userPostfix, updateUserHandler)
	r.Delete(userPostfix, deleteUserHandler)

	server := http.Server{
		Addr:         baseURL,
		Handler:      r,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found")
	}
}
