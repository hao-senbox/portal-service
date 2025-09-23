package bmi

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	
	if date != nil {
		var bmi BMI
		opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
		err := b.collection.FindOne(ctx, filter, opts).Decode(&bmi)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return []*BMI{}, nil
			}
			return nil, err
		}
		return []*BMI{&bmi}, nil
	}

	cursor, err := b.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bmis []*BMI
	if err := cursor.All(ctx, &bmis); err != nil {
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