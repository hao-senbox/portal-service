package selectoptions

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SelectOptionsRepository interface {
	Create(ctx context.Context, doc *SelectOptions) error
	GetSelectOption(ctx context.Context, typeStr, organizationID, termID, studentID string) (*SelectOptions, error)
}

type selectOptionsRepository struct {
	collection *mongo.Collection
}

func NewSelectOptionsRepository(collection *mongo.Collection) SelectOptionsRepository {
	return &selectOptionsRepository{
		collection: collection,
	}
}

func (r *selectOptionsRepository) Create(ctx context.Context, doc *SelectOptions) error {
	filter := bson.M{
		"organization_id": doc.OrganizationID,
		"term_id":         doc.TermID,
		"student_id":      doc.StudentID,
		"type":            doc.Type,
	}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, filter, bson.M{"$set": doc}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (r *selectOptionsRepository) GetSelectOption(ctx context.Context, typeStr, organizationID, termID, studentID string) (*SelectOptions, error) {
	var selectOption *SelectOptions
	fmt.Printf("typeStr: %s, organizationID: %s, termID: %s, studentID: %s\n", typeStr, organizationID, termID, studentID)
	filter := bson.M{
		"organization_id": organizationID,
		"term_id":         termID,
		"student_id":      studentID,
		"type":            typeStr,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&selectOption)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return selectOption, nil
}