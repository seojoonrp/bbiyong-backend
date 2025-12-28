package database

import (
	"context"
	"log"
	"time"

	"github.com/seojoonrp/bbiyong-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Successfully Connected to MongoDB!")

	// createGeoIndex(client)
	createUserIndex(client)

	return client, nil
}

// 2dsphere 인덱스 생성
func createGeoIndex(client *mongo.Client) {
	coll := client.Database(config.AppConfig.DBName).Collection("meetings")

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: "2dsphere"}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Could not create 2dsphere index: %v", err)
	} else {
		log.Println("2dsphere index ensured on meetings collection")
	}
}

func createUserIndex(client *mongo.Client) {
	coll := client.Database(config.AppConfig.DBName).Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Could not create unique index on username: %v", err)
	} else {
		log.Println("Unique index ensured on username field in users collection")
	}
}
