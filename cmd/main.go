package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/go-chi/chi"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURL      = "localhost:8081"
	usersPostfix = "/users"
	userPostfix  = usersPostfix + "/{id}"
)

// Role is ...
type Role uint8

const (
	// AdminRole is ...
	AdminRole Role = iota
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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := &NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	id := createUser(newUser)

	fmt.Println("new user id:", id)
}

func createUser(user *NewUserData) int64 {
	id := generateUserID()

	fmt.Printf("new user data: %+v\n", *user)

	return id
}

func generateUserID() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))
	if err != nil {
		panic(err)
	}

	return nBig.Int64()
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	user := getUser(userID)

	fmt.Printf("get user: %+v\n", *user)
}

func getUser(id int64) *UserData {
	fmt.Printf("get user id: %v\n", id)

	return &UserData{
		ID:        id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
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

	server := http.Server{
		Addr:         baseURL,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
