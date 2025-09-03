package body

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BodyRepository interface {
	PushCheckIn(ctx context.Context, checkIn *CheckIn) error
	GetCheckIns(ctx context.Context, student_id string, date *time.Time) ([]*CheckIn, error)
}

type bodyRepository struct {
	collection *mongo.Collection
}

func NewBodyRepository(collection *mongo.Collection) BodyRepository {
	return &bodyRepository{
		collection: collection,
	}
}

func (r *bodyRepository) PushCheckIn(ctx context.Context, checkIn *CheckIn) error {

	filter := bson.M{
		"student_id": checkIn.StudentID,
		"type":       checkIn.Type,
		"date":       checkIn.Date,
		"context":    checkIn.Context,
		"gender":     checkIn.Gender,
	}

	update := bson.M{
		"$setOnInsert": bson.M{
			"created_by": checkIn.CreatedBy,
			"created_at": time.Now(),
		},
		"$push": bson.M{
			"marks": bson.M{
				"$each": checkIn.Marks,
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err

}

func (r *bodyRepository) GetCheckIns(ctx context.Context, student_id string, date *time.Time) ([]*CheckIn, error) {

	filter := bson.M{
		"student_id": student_id,
	}

	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		filter["date"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}

	var checkIns []*CheckIn
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &checkIns)
	if err != nil {
		return nil, err
	}

	return checkIns, nil
	
}