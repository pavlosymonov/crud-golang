package main

import (
	"fmt"
	"main/internal/http"
	"main/internal/services/db"
	grpcclient "main/internal/services/grpc-client"
	grpcserver "main/internal/services/grpc-server"
	"main/internal/services/user"
)

func main() {
	fmt.Println("Starting the application...")

	// Services
	db := db.New().Database("store-nest")
	userService := user.New(db)

	// Run grpc server in go routine
	go grpcserver.New(userService)

	grpcUserClient := grpcclient.New()

	// Server
	server := http.New(grpcUserClient)
	server.SetupRouts()
	server.Run()
}
