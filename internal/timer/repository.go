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
	CreateIsTime(ctx context.Context, isTime *IsTime) error
	GetIsTimes(ctx context.Context, studentID string) ([]*IsTime, error)
}

type timerRepository struct {
	collection *mongo.Collection
	IsTimeCollection *mongo.Collection
}

func NewTimerRepository(collection *mongo.Collection, IsTimeCollection *mongo.Collection) TimerRepository {
	return &timerRepository{
		collection: collection,
		IsTimeCollection: IsTimeCollection,
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

func (t *timerRepository) CreateIsTime(ctx context.Context, isTime *IsTime) error {

	_, err := t.IsTimeCollection.InsertOne(ctx, isTime)
	if err != nil {
		return err
	}

	return nil

}

func (t *timerRepository) GetIsTimes(ctx context.Context, studentID string) ([]*IsTime, error) {

	filter := bson.M{}

	if studentID != "" {
		filter["student_id"] = studentID
	}

	findOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var isTimes []*IsTime

	cursor, err := t.IsTimeCollection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &isTimes)
	if err != nil {
		return nil, err
	}

	return isTimes, nil
	
}