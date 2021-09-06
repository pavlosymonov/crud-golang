package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"main/internal/models/user"
)

var (
	ErrInternal = errors.New("internal error")
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)

type User interface {
	Create(ctx context.Context, c user.User) (*mongo.InsertOneResult, error)
	List(ctx context.Context) ([]user.User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type Service struct {
	db *mongo.Client
}

func New(db *mongo.Client) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Create(ctx context.Context, user user.User) (*mongo.InsertOneResult, error) {
	collection := s.db.Database("store-nest").Collection("users")

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) List(ctx context.Context) ([]user.User, error) {
	var users = make([]user.User, 0)
	collection := s.db.Database("store-nest").Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user user.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.Database("store-nest").Collection("users")

	qry := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, qry)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
