package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"user-service/internal/api/grpc_server"
	"user-service/internal/api/handlers"
	"user-service/internal/api/routes"
	"user-service/internal/services"
	"user-service/internal/storage/postgres"
	pb "user-service/proto"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

    connectionString := os.Getenv("DATABASE_URL")
    db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Connection to the database failed: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserStorage(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := mux.NewRouter()
    routes.RegisterUserHandlers(router, userHandler)

    // RunHttpServer(router)
	RunGrpcServer(grpc_server.NewServer(userService))
}

func RunHttpServer(router *mux.Router) {
	httpServer := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    log.Println("Server run :8080...")
    log.Fatal(httpServer.ListenAndServe())
}

func RunGrpcServer(server pb.UserserviceServer) {
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserserviceServer(grpcServer, server)

	log.Println("gRPC Server run :5000...")
	grpcServer.Serve(lis)
}
