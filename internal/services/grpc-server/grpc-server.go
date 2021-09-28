package grpcserver

import (
	"google.golang.org/grpc"
	"log"
	"main/internal/services/user"
	userpb "main/proto"
	"net"
	"os"
	"os/signal"
)

func New(userService *user.Service) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, userService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC server on port 9000: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	grpcServer.Stop()
	lis.Close()
}
