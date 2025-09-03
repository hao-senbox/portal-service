package body

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type BodyRepository interface {
	CreateCheckIn(ctx context.Context, checkIn *CheckIn) (string, error)
}

type bodyRepository struct {
	collection *mongo.Collection
}

func NewBodyRepository(collection *mongo.Collection) BodyRepository {
	return &bodyRepository{
		collection: collection,
	}
}

func (r *bodyRepository) CreateCheckIn(ctx context.Context, checkIn *CheckIn) (string, error) {
	return "", nil
}
