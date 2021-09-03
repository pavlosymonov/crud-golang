package http

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	userpb "main/proto"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	router *mux.Router
	userClient userpb.UserServiceClient
}

func New(grpcUserClient userpb.UserServiceClient) *Server {
	return &Server{
		router: mux.NewRouter(),
		userClient: grpcUserClient,
	}
}

func (s *Server) SetupRouts() {
	apiRouter := s.router.PathPrefix("/api/v0").Subrouter()

	apiRouter.HandleFunc("/user", s.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users", s.GetUsers).Methods("GET")
	apiRouter.HandleFunc("/user/{id}", s.DeleteUser).Methods("DELETE")
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