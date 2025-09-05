package timer

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimerRepository interface {
	CreateTimer(ctx context.Context, timer *Timer) (string, error)
	GetTimers(ctx context.Context, studentID string) ([]*Timer, error)
}

type timerRepository struct {
	collection *mongo.Collection
}

func NewTimerRepository(collection *mongo.Collection) TimerRepository {
	return &timerRepository{
		collection: collection,
	}
}

func (t *timerRepository) CreateTimer(ctx context.Context, timer *Timer) (string, error) {

	result, err := t.collection.InsertOne(ctx, timer)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (t *timerRepository) GetTimers(ctx context.Context, studentID string) ([]*Timer, error) {

	filter := bson.M{}

	if studentID != "" {
		filter["student_id"] = studentID
	}

	findOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var timers []*Timer

	cursor, err := t.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &timers)
	if err != nil {
		return nil, err
	}

	return timers, nil
}
