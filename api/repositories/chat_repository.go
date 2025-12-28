// api/repositories/chat_repository.go

package repositories

import (
	"context"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *models.ChatMessage) error
	GetChatHistory(ctx context.Context, meetingID primitive.ObjectID, limit int64) ([]models.ChatMessage, error)
}

type chatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) ChatRepository {
	return &chatRepository{
		collection: db.Collection("chats"),
	}
}

func (r *chatRepository) SaveMessage(ctx context.Context, msg *models.ChatMessage) error {
	_, err := r.collection.InsertOne(ctx, msg)
	return err
}

func (r *chatRepository) GetChatHistory(ctx context.Context, meetingID primitive.ObjectID, limit int64) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	filter := bson.M{"meeting_id": meetingID}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
