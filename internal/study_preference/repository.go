package studypreference

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudyPreferenceRepository interface {
	CreateStudyPreference(ctx context.Context, studyPreference *StudyPreference) error
	GetStudyPreferencesByStudentID(ctx context.Context, studentID, termID, orgID string) (*StudyPreference, error)
	UpdateStudyPreference(ctx context.Context, id primitive.ObjectID, updateData map[string]interface{}) error
	GetStudyPreferenceByID(ctx context.Context, id primitive.ObjectID) (*StudyPreference, error)
}

type studyPreferenceRepository struct {
	collection *mongo.Collection
}

func NewStudyPreferenceRepository(collection *mongo.Collection) StudyPreferenceRepository {
	return &studyPreferenceRepository{collection: collection}
}

func (r *studyPreferenceRepository) CreateStudyPreference(ctx context.Context, studyPreference *StudyPreference) error {
	
	filter := bson.M{
		"student_id": studyPreference.StudentID,
		"term_id": studyPreference.TermID,
		"organization_id": studyPreference.OrganizationID,
		"is_deleted": false,
	}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	result, err := r.collection.InsertOne(ctx, studyPreference)
	if err != nil {
		return err
	}

	studyPreference.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *studyPreferenceRepository) GetStudyPreferencesByStudentID(ctx context.Context, studentID, termID, orgID string) (*StudyPreference, error) {
	var studyPreference *StudyPreference
	
	filter := bson.M{
		"student_id": studentID,
		"term_id": termID,
		"is_deleted": false,
		"organization_id": orgID,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&studyPreference)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return studyPreference, nil
}

func (r *studyPreferenceRepository) UpdateStudyPreference(ctx context.Context, id primitive.ObjectID, updateData map[string]interface{}) error {
	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{"$set": updateData}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *studyPreferenceRepository) GetStudyPreferenceByID(ctx context.Context, id primitive.ObjectID) (*StudyPreference, error) {

	var studyPreference StudyPreference

	filter := bson.M{"_id": id, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&studyPreference)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("study preference not found")
		}
		return nil, err
	}

	return &studyPreference, nil
}
