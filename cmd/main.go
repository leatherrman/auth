package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/go-chi/chi"
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
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateUserData is ...
type UpdateUserData struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

// ResponseUserID is ...
type ResponseUserID struct {
	ID int64 `json:"id"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := &NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	id, errId := createUser(newUser)

	if errId != nil {
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

func createUser(user *NewUserData) (int64, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		return -1, err
	}

	fmt.Printf("new user data: %+v\n", *user)

	return nBig.Int64(), nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := getUser(userID)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode new user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getUser(id int64) *UserData {
	fmt.Printf("get user id: %v\n", id)

	return &UserData{
		ID:        id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}
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

	deleteUser(userID)

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(id int64) {
	fmt.Printf("delete user id: %v\n", id)
}

func parseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func main() {
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
