package bmi

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BMIRepo interface {
	CreateBMI(ctx context.Context, bmi *BMI) (string, error)
	GetBMIs(ctx context.Context, student_id string, date *time.Time) ([]*BMI, error)
	GetBMI(ctx context.Context, id primitive.ObjectID) (*BMI, error)
}

type bmiRepository struct {
	collection *mongo.Collection
}

func NewBMIRepository(collection *mongo.Collection) BMIRepo {
	return &bmiRepository{
		collection: collection,
	}
}

func (b *bmiRepository) CreateBMI(ctx context.Context, bmi *BMI) (string, error) {

	result, err := b.collection.InsertOne(ctx, bmi)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
	
}

func (b *bmiRepository) GetBMIs(ctx context.Context, student_id string, date *time.Time) ([]*BMI, error) {

	filter := bson.M{}

	if student_id != "" {
		filter["student_id"] = student_id
	}

	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		filter["date"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}

	var bmis []*BMI

	cursor, err := b.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &bmis)
	if err != nil {
		return nil, err
	}

	return bmis, nil
	
}

func (b *bmiRepository) GetBMI(ctx context.Context, id primitive.ObjectID) (*BMI, error) {

	var bmi BMI

	err := b.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bmi)
	if err != nil {
		return nil, err
	}

	return &bmi, nil
	
}