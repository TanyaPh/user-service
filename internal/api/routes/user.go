package routes

import (
	"github.com/gorilla/mux"

	"user-service/internal/api/handlers"
)

func RegisterUserHandlers(router *mux.Router, handler *handlers.User) {
    router.HandleFunc("/api/users", handler.CreateUserHandler).Methods("POST")
    router.HandleFunc("/api/users/{id}", handler.GetUserByIdHandler).Methods("GET")
    router.HandleFunc("/api/users", handler.GetAllUsersHandler).Methods("GET")
    router.HandleFunc("/api/users/{id}", handler.EditUserHandler).Methods("PUT")
    router.HandleFunc("/api/users/{id}", handler.DeleteUserHandler).Methods("DELETE")
}
