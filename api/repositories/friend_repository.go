// api/repositories/friend_repository.go

package repositories

import (
	"context"
	"time"

	"github.com/seojoonrp/bbiyong-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendRepository interface {
	SendRequest(ctx context.Context, f *models.Friendship) error
	FindByID(ctx context.Context, fID primitive.ObjectID) (*models.Friendship, error)
	FindByUserIDs(ctx context.Context, uID1, uID2 primitive.ObjectID) (*models.Friendship, error)
	UpdateStatus(ctx context.Context, fID primitive.ObjectID, status string) error
	GetFriends(ctx context.Context, uID primitive.ObjectID) ([]models.Friendship, error)
}

type friendRepository struct {
	collection *mongo.Collection
}

func NewFriendRepository(db *mongo.Database) FriendRepository {
	return &friendRepository{collection: db.Collection("friendships")}
}

func (r *friendRepository) SendRequest(ctx context.Context, f *models.Friendship) error {
	_, err := r.collection.InsertOne(ctx, f)
	return err
}

func (r *friendRepository) FindByID(ctx context.Context, fID primitive.ObjectID) (*models.Friendship, error) {
	var f models.Friendship
	err := r.collection.FindOne(ctx, bson.M{"_id": fID}).Decode(&f)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *friendRepository) FindByUserIDs(ctx context.Context, uID1, uID2 primitive.ObjectID) (*models.Friendship, error) {
	var f models.Friendship
	filter := bson.M{
		"$or": []bson.M{
			{"requester_id": uID1, "addressee_id": uID2},
			{"requester_id": uID2, "addressee_id": uID1},
		},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&f)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &f, err
}

func (r *friendRepository) UpdateStatus(ctx context.Context, fID primitive.ObjectID, status string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": fID}, bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	return err
}

func (r *friendRepository) GetFriends(ctx context.Context, uID primitive.ObjectID) ([]models.Friendship, error) {
	var results []models.Friendship
	filter := bson.M{
		"status": models.FriendStatusAccepted,
		"$or":    []bson.M{{"requester_id": uID}, {"addressee_id": uID}},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	return results, err
}
