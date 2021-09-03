package grpcclient

import (
	"google.golang.org/grpc"
	"log"
	userpb "main/proto"
)

func New() userpb.UserServiceClient {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect: %v", err)
	}

	c := userpb.NewUserServiceClient(conn)

	return c
}
