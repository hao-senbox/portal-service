package body

import (
	"context"
	"fmt"
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

	now := time.Now()

	for _, mark := range checkIn.Marks {

		markData := bson.M{
			"name":         mark.Name,
			"note":         mark.Note,
			"color":        mark.Color,
			"severity":     mark.Severity,
			"submitted_at": now,
		}

		filterWithMark := bson.M{
			"student_id": checkIn.StudentID,
			"type":       checkIn.Type,
			"date":       checkIn.Date,
			"context":    checkIn.Context,
			"gender":     checkIn.Gender,
			"marks.name": mark.Name, 
		}

		updateExisting := bson.M{
			"$set": bson.M{
				"marks.$":    markData, 
				"updated_at": now,
			},
		}

		result, err := r.collection.UpdateOne(ctx, filterWithMark, updateExisting)
		if err != nil {
			return fmt.Errorf("error when update mark: %v", err)
		}

		if result.MatchedCount == 0 {
			ensureDoc := bson.M{
				"$setOnInsert": bson.M{
					"student_id": checkIn.StudentID,
					"type":       checkIn.Type,
					"date":       checkIn.Date,
					"context":    checkIn.Context,
					"gender":     checkIn.Gender,
					"created_by": checkIn.CreatedBy,
					"created_at": now,
					"marks":      []interface{}{}, 
				},
				"$set": bson.M{
					"updated_at": now,
				},
			}

			_, err = r.collection.UpdateOne(ctx, filter, ensureDoc, options.Update().SetUpsert(true))
			if err != nil {
				return fmt.Errorf("error when ensure doc: %v", err)
			}

			addMark := bson.M{
				"$push": bson.M{
					"marks": markData,
				},
				"$set": bson.M{
					"updated_at": now,
				},
			}

			_, err = r.collection.UpdateOne(ctx, filter, addMark)
			if err != nil {
				return fmt.Errorf("error when add mark: %v", err)
			}
		}
	}

	return nil
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