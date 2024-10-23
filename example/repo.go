package main

import (
	"context"
	"github.com/its-own/gaudit/db"
	"github.com/its-own/gaudit/in"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	in.Inject
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

// IUserRepo is a User repository
type IUserRepo interface {
	Create(ctx context.Context, param *User) (*User, error)
}

// UserRepo implementation of IUserRepo, also holds collection name and mongo db rapper repository
type UserRepo struct {
	collection string
	connection db.NoSql
}

func NewUserRepo(collection string, connection db.NoSql) IUserRepo {
	return &UserRepo{connection: connection, collection: collection}
}

// Create is a simple implementation of user repository
func (u UserRepo) Create(ctx context.Context, param *User) (*User, error) {
	err := u.connection.Insert(ctx, u.collection, param)
	if err != nil {
		return nil, err
	}
	return param, nil
}
