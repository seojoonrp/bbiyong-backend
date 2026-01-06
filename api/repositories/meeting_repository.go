// api/repositories/meeting_repository.go

package repositories

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting *models.Meeting) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Meeting, error)
	FindNearby(ctx context.Context, lon, lat float64, radiusMeter float64, days []int) ([]models.Meeting, error)
	AddParticipant(ctx context.Context, meetingID, userID primitive.ObjectID, maxParticipants int) (bool, error)
	RemoveParticipant(ctx context.Context, meetingID, userID primitive.ObjectID, maxParticipants int) (bool, error)
	IncrementSaveCount(ctx context.Context, meetingID primitive.ObjectID) error
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

func (r *meetingRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Meeting, error) {
	var meeting models.Meeting
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meeting)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &meeting, nil
}

func (r *meetingRepository) FindNearby(ctx context.Context, lon, lat float64, radiusMeter float64, days []int) ([]models.Meeting, error) {
	var meetings []models.Meeting

	// 몽고디비의 개쩌는 공간 쿼리
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

	if len(days) > 0 {
		filter["day_of_week"] = bson.M{"$in": days}
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

func (r *meetingRepository) AddParticipant(ctx context.Context, meetingID primitive.ObjectID, userID primitive.ObjectID, maxParticipants int) (bool, error) {
	filter := bson.M{
		"_id":    meetingID,
		"status": models.MeetingStatusRecruiting,
		"participants." + strconv.Itoa(maxParticipants-1): bson.M{"$exists": false}, // 마지막 원소가 있는지 확인 -> 정원 초과 여부를 확인할 수 있음
		"participants": bson.M{"$ne": userID},
	}
	update := bson.M{"$addToSet": bson.M{"participants": userID}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}

	fullFilter := bson.M{
		"_id":    meetingID,
		"status": models.MeetingStatusRecruiting,
		"participants." + strconv.Itoa(maxParticipants-1): bson.M{"$exists": true},
	}
	fullUpdate := bson.M{"$set": bson.M{"status": models.MeetingStatusFull}}
	_, err = r.collection.UpdateOne(ctx, fullFilter, fullUpdate)
	if err != nil {
		log.Println("Successfully added user to meeting, but error occurred while updating status to full:", err)
	}

	return true, nil
}

func (r *meetingRepository) RemoveParticipant(ctx context.Context, meetingID primitive.ObjectID, userID primitive.ObjectID, maxParticipants int) (bool, error) {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": meetingID},
		bson.M{"$pull": bson.M{"participants": userID}},
	)
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}

	backFilter := bson.M{
		"_id":    meetingID,
		"status": models.MeetingStatusFull,
		"participants." + strconv.Itoa(maxParticipants-1): bson.M{"$exists": false},
	}
	backUpdate := bson.M{"$set": bson.M{"status": models.MeetingStatusRecruiting}}

	_, err = r.collection.UpdateOne(ctx, backFilter, backUpdate)
	if err != nil {
		log.Println("Successfully removed user from meeting, but error occurred while updating status to recruiting:", err)
	}

	return true, nil
}

func (r *meetingRepository) IncrementSaveCount(ctx context.Context, meetingID primitive.ObjectID) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": meetingID},
		bson.M{"$inc": bson.M{"save_count": 1}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("meeting not found")
	}

	return nil
}
