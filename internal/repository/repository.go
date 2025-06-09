package repository

import (
	"context"
	"portal/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PortalRepository interface {
	CreateStudentActivity(ctx context.Context, activityStudent *models.StudentActivity) error	
	GetAllStudentActivity(ctx context.Context, studentID string, date *time.Time) ([]*models.StudentActivity, error)
}

type portalRepository struct {
	collection *mongo.Collection
}

func NewPortalRepository(collection *mongo.Collection) PortalRepository {
	return &portalRepository {
		collection: collection,
	}
}

func (r *portalRepository) CreateStudentActivity(ctx context.Context, activityStudent *models.StudentActivity) error {

	_, err := r.collection.InsertOne(ctx, activityStudent)
	if err != nil {
		return err
	}

	return nil
}

func (r *portalRepository) GetAllStudentActivity(ctx context.Context, studentID string, date *time.Time) ([]*models.StudentActivity, error) {

	var activities []*models.StudentActivity

	filter := bson.M{"student_id": studentID}

	if date != nil {
		filter["date"] = date
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}

	return activities, nil

}
