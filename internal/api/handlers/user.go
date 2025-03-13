package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/gorilla/mux"
)

type User struct {
	service *services.User
}

func NewUserHandler(userService *services.User) *User {
	return &User{
		service: userService,
	}
}

func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
	}

	var request requestBody
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := models.NewUser(request.Name, request.Address, request.Phone)
	userId, err := u.service.CreateUser(user)
	if err != nil {
		log.Println("Couldn't create user. Error: %w", err)
		http.Error(w, "Can't create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"id": userId})
	w.WriteHeader(http.StatusOK)
}

func (u *User) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIdStr := params["id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := u.service.GetUserById(userId)
	if err != nil {
		log.Println("Couldn't get user. Error: %w", err)
		http.Error(w, "Can't get user ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.GetUsers()
	if err != nil {
		log.Println("Couldn't get users. Error: %w", err)
		http.Error(w, "Cann't get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *User) EditUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIdStr := params["id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	type requestBody struct {
		Name    *string `json:"name"`
		Address *string `json:"address"`
		Phone   *string `json:"phone"`
	}

	var request requestBody
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = u.service.EditUser(userId, request.Name, request.Address, request.Phone)
	if err != nil {
		log.Println("Couldn't edit user. Error: %w", err)
		http.Error(w, "Cann't edit user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *User) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIdStr := params["id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = u.service.DeleteUser(userId)
	if err != nil {
		log.Println("Couldn't delete user. Error: %w", err)
		http.Error(w, "Cann't delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
