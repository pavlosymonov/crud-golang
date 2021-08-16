package main

import (
	"fmt"
	"main/internal/http"
	"main/internal/services/db"
	"main/internal/services/user"
)

func main() {
	fmt.Println("Starting the application...")

	// Services
	client := db.New()
	userService := user.New(client)

	// Server
	server := http.New(userService)

	server.SetupRouts()
	server.Run()
}
