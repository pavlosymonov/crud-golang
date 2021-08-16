package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main/internal/models/user"
)

var (
	ErrAlreadyExist = errors.New("user already exist")
	ErrNotExist     = errors.New("user doesn't exist")
)

type User interface {
	Create(ctx context.Context, c user.User) (*mongo.InsertOneResult, error)
	List(ctx context.Context) ([]user.User, error)
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
	var users []user.User
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


//func (service *Service) Update(ctx context.Context, id int, c user.User) (*user.User, error) {
//	res, err := service.db.UpdateCounty(ctx, id, c)
//	if err == storage.ErrAlreadyExist {
//		return nil, ErrAlreadyExist
//	}
//	if err == storage.ErrNotExist {
//		return nil, ErrNotExist
//	}
//	if err != nil {
//		return nil, fmt.Errorf("county updating failed: %w", err)
//	}
//
//	return res, nil
//}

//func (service *Service) Delete(ctx context.Context, id int) error {
//	err := service.s.DeleteCounty(ctx, id)
//	if err == storage.ErrNotExist {
//		return ErrNotExist
//	}
//	if err != nil {
//		return fmt.Errorf("county deleting failed: %w", err)
//	}
//
//	return nil
//}
//
//func (service *Service) Read(ctx context.Context, id int) (*county.County, error) {
//	res, err := service.s.ReadCounty(ctx, id)
//	if err == storage.ErrNotExist {
//		return nil, ErrNotExist
//	}
//	if err != nil {
//		return nil, fmt.Errorf("county reading failed: %w", err)
//	}
//
//	return res, nil
//}
//
