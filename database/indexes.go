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

func EnsureIndexes(client *mongo.Client) {
	db := client.Database(config.AppConfig.DBName)

	initUserIndexes(db.Collection("users"))
	initMeetingIndexes(db.Collection("meetings"))
	initChatIndexes(db.Collection("chats"))
	initFriendshipIndexes(db.Collection("friendships"))
}

func initUserIndexes(coll *mongo.Collection) {
	// Unique username
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_unique_username"),
	}
	createIndex(coll, indexModel)
}

func initMeetingIndexes(coll *mongo.Collection) {
	// Geospatial location
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
		Options: options.Index().SetName("idx_geo_location"),
	}
	createIndex(coll, indexModel)
}

func initChatIndexes(coll *mongo.Collection) {
	// 모임별 최신 채팅 조회
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "meeting_id", Value: 1},
			{Key: "created_at", Value: -1},
		},
		Options: options.Index().SetName("idx_meeting_id_created_at"),
	}
	createIndex(coll, indexModel)
}

func initFriendshipIndexes(coll *mongo.Collection) {
	// Unique request-address pair
	createIndex(coll, mongo.IndexModel{
		Keys: bson.D{
			{Key: "requester_id", Value: 1},
			{Key: "addressee_id", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("idx_unique_friendship"),
	})
	// 나한테 온 요청 조회
	createIndex(coll, mongo.IndexModel{
		Keys: bson.D{
			{Key: "addressee_id", Value: 1},
			{Key: "status", Value: 1},
		},
		Options: options.Index().SetName("idx_received_requests"),
	})
}

func createIndex(coll *mongo.Collection, model mongo.IndexModel) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	name, err := coll.Indexes().CreateOne(ctx, model)
	if err != nil {
		log.Printf("Error while creating index on %s: %v", coll.Name(), err)
		return
	}
	log.Printf("Successfully applied index %s on collection %s", name, coll.Name())
}
