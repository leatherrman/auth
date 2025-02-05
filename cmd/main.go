package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/go-chi/chi"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	baseUrl      = "localhost:8081"
	usersPostfix = "/users"
	userPostfix  = usersPostfix + "/{id}"
)

const (
	admin uint8 = iota
	user
)

type NewUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            uint8  `json:"role"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      uint8     `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := &NewUser{}
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	id := createUser(newUser)

	fmt.Println("new user id:", id)
}

func createUser(user *NewUser) int64 {
	id := rand.Int63()

	fmt.Printf("new user data: %+v\n", *user)

	return id
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

func getUser(id int64) *User {
	fmt.Printf("get user id: %v\n", id)

	return &User{
		ID:        id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      user,
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

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
