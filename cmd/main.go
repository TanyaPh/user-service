package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	var wg sync.WaitGroup

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signalChan
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunHttpServer(ctx, router)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunGrpcServer(ctx, grpc_server.NewServer(userService))
	}()

	wg.Wait()
}

func RunHttpServer(ctx context.Context, router *mux.Router) {
	httpServer := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP Server...")
		httpServer.Shutdown(context.Background())
	}()

    log.Println("Server run :8080...")
    log.Fatal(httpServer.ListenAndServe())
}

func RunGrpcServer(ctx context.Context, server pb.UserserviceServer) {
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserserviceServer(grpcServer, server)

	go func() {
		<-ctx.Done()
		log.Println("Shutting down gRPC Server...")
		grpcServer.GracefulStop()
	}()

	log.Println("gRPC Server run :5000...")
	grpcServer.Serve(lis)
}
