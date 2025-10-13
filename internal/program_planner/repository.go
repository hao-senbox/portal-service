package program_planner

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProgramPlanerRepository interface {
	CreateProgramPlaner(ctx context.Context, data *ProgramPlaner) (string, error)
	GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error)
	GetProgramPlaner(ctx context.Context, id primitive.ObjectID) (*ProgramPlaner, error)
	UpdateProgramPlaner(ctx context.Context, data *ProgramPlaner, id primitive.ObjectID) error
	DeleteProgramPlaner(ctx context.Context, id primitive.ObjectID) error
}

type programPlannerRepository struct {
	programPlanerCollection *mongo.Collection
}

func NewProgramPlanerRepository(collection *mongo.Collection) ProgramPlanerRepository {
	return &programPlannerRepository{
		programPlanerCollection: collection,
	}
}

func (repository *programPlannerRepository) CreateProgramPlaner(ctx context.Context, data *ProgramPlaner) (string, error) {

	result, err := repository.programPlanerCollection.InsertOne(ctx, data)

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", err
	}

	return oid.Hex(), err
}

func (repository *programPlannerRepository) GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error) {

	var programPlaner []*ProgramPlaner

	filter := bson.M{
		"is_deleted": false,
	}

	cursor, err := repository.programPlanerCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &programPlaner)

	return programPlaner, err

}

func (repository *programPlannerRepository) GetProgramPlaner(ctx context.Context, id primitive.ObjectID) (*ProgramPlaner, error) {

	var programPlaner ProgramPlaner

	filter := bson.M{
		"_id": id,
	}

	err := repository.programPlanerCollection.FindOne(ctx, filter).Decode(&programPlaner)

	if err != nil {
		return nil, err
	}

	return &programPlaner, nil

}

func (repository *programPlannerRepository) UpdateProgramPlaner(ctx context.Context, data *ProgramPlaner, id primitive.ObjectID) error {
	
	filter := bson.M{"_id": id}
	update := bson.M{"$set": data}
	
	_, err := repository.programPlanerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	return nil

}

func (repository *programPlannerRepository) DeleteProgramPlaner(ctx context.Context, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_deleted": true}}
	
	_, err := repository.programPlanerCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	return nil
}