package main

import (
	"context"
	"fmt"
	"github.com/its-own/gaudit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

func main() {
	// connect to mongo db
	ctx := context.Background()
	client := connectMongo(ctx, "mongodb://localhost:27017")
	// initialize go audit
	aMgo := gaudit.Init(&gaudit.Config{
		Client:   client,
		Database: client.Database("test_database"),
		Logger:   slog.Default(),
	})
	// create user and pass gaudit mongo instance
	_, err := NewUserRepo("user", aMgo).Create(ctx, &User{
		ID:   primitive.NewObjectID(),
		Name: "Razibul Hasan Mithu",
	})
	if err != nil {
		return
	}
}

// connectMongo connect to mongo db
func connectMongo(ctx context.Context, mongoUrl string) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl).SetRetryWrites(false))
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	slog.Default().Info("Connected to MongoDB")
	return client
}
