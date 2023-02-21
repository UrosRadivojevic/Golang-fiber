package user_repository

import (
	"context"
	"fmt"

	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Interface interface {
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Register(ctx context.Context, user model.User) error
	Exists(ctx context.Context, username string) bool
}

type Repository struct {
	col *mongo.Collection
}

func New(col *mongo.Collection) *Repository {
	return &Repository{
		col: col,
	}
}

func (u *Repository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	filter := bson.M{"username": username}
	val := u.col.FindOne(ctx, filter)
	var user model.User
	if err := val.Decode(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *Repository) Register(ctx context.Context, user model.User) error {
	exist := u.Exists(ctx, user.Username)
	if exist {
		fmt.Println("Username already exists!")
		return nil
	}
	_, err := u.col.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func (u *Repository) Exists(ctx context.Context, username string) bool {
	filter := bson.M{"username": username}
	val := u.col.FindOne(ctx, filter)
	var user model.User
	if err := val.Decode(&user); err != nil {
		return false
	}
	return true
}
