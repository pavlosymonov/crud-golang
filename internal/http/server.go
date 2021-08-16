package http

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"main/internal/services/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	router             *mux.Router
	userService        user.User
}

func New(userService user.User) *Server {
	return &Server{
		router: mux.NewRouter(),
		userService: userService,
	}
}

func (s *Server) SetupRouts() {
	apiRouter := s.router.PathPrefix("/api/v0").Subrouter()

	apiRouter.HandleFunc("/user", s.CreateUserEndpoint).Methods("POST")
	apiRouter.HandleFunc("/users", s.GetUsersEndpoint).Methods("GET")
}


func (s *Server) Run() {
	server := &http.Server{
		Addr:         ":8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      s.router,
	}

	// Graceful shutdown.

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Launching on address: 8081")
		log.Println(server.ListenAndServe())
		signalChan <- syscall.SIGTERM
	}()

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down")
	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}