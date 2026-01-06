// api/repositories/save_repository.go

package repositories

import (
	"context"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaveRepository interface {
	Create(ctx context.Context, save *models.Save) error
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
