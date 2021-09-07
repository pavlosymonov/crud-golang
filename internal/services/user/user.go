package user

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"main/internal/models/user"
	userpb "main/proto"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)

type User interface {
	CreateUser(context.Context, *userpb.CreateUserReq) (*userpb.CreateUserRes, error)
	DeleteUser(context.Context, *userpb.DeleteUserReq) (*userpb.DeleteUserRes, error)
	ListUsers(*userpb.ListUsersReq, userpb.UserService_ListUsersServer) error
}

type Service struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateUser(ctx context.Context, req *userpb.CreateUserReq) (*userpb.CreateUserRes, error) {
	collection := s.db.Collection("users")

	userproto := req.GetUser()

	data := user.User{
		Login: userproto.GetLogin(),
		Email:    userproto.GetEmail(),
		BillingAddress:  userproto.GetBillingAddress(),
		ShippingAddress:  userproto.GetShippingAddress(),
		Phone:  userproto.GetPhone(),
	}

	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	oid := result.InsertedID.(primitive.ObjectID)
	userproto.Id = oid.Hex()

	return &userpb.CreateUserRes{User: userproto}, nil
}

func (s *Service) ListUsers(req *userpb.ListUsersReq, stream userpb.UserService_ListUsersServer) error {
	data := &user.User{}
	collection := s.db.Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		stream.Send(&userpb.ListUsersRes{
			User: &userpb.User{
				Id:       data.ID.Hex(),
				Login: data.Login,
				Email:  data.Email,
				BillingAddress:    data.BillingAddress,
				ShippingAddress: data.ShippingAddress,
				Phone: data.Phone,
			},
		})

	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, req *userpb.DeleteUserReq) (*userpb.DeleteUserRes, error) {
	collection := s.db.Collection("users")

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	qry := bson.M{"_id": oid}

	_, err = collection.DeleteOne(ctx, qry)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}

	return &userpb.DeleteUserRes{
		Success: true,
	}, nil
}
