// api/repositories/save_repository.go

package repositories

import (
	"context"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaveRepository interface {
	Create(ctx context.Context, save *models.Save) error
	Delete(ctx context.Context, userID, meetingID primitive.ObjectID) (int64, error)
}

type saveRepository struct {
	collection *mongo.Collection
}

func NewSaveRepository(db *mongo.Database) SaveRepository {
	return &saveRepository{collection: db.Collection("saves")}
}

func (r *saveRepository) Create(ctx context.Context, save *models.Save) error {
	_, err := r.collection.InsertOne(ctx, save)
	return err
}

func (r *saveRepository) Delete(ctx context.Context, userID, meetingID primitive.ObjectID) (int64, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"user_id":    userID,
		"meeting_id": meetingID,
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
