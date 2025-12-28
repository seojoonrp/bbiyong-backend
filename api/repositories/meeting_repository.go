// api/repositories/meeting_repository.go

package repositories

import (
	"context"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting *models.Meeting) error
	FindNearby(ctx context.Context, lon, lat float64, radiusMeter float64) ([]models.Meeting, error)
}

type meetingRepository struct {
	collection *mongo.Collection
}

func NewMeetingRepository(db *mongo.Database) MeetingRepository {
	return &meetingRepository{
		collection: db.Collection("meetings"),
	}
}

func (r *meetingRepository) Create(ctx context.Context, meeting *models.Meeting) error {
	_, err := r.collection.InsertOne(ctx, meeting)
	return err
}

func (r *meetingRepository) FindNearby(ctx context.Context, lon, lat float64, radiusMeter float64) ([]models.Meeting, error) {
	var meetings []models.Meeting

	// MongoDB의 개쩌는 공간 쿼리
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat},
				},
				"$maxDistance": radiusMeter,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &meetings); err != nil {
		return nil, err
	}
	return meetings, nil
}
