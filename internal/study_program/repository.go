package studyprogram

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudyProgramRepository interface {
	CreateStudyProgram(ctx context.Context, data *StudyProgram) error
	GetStudyPrograms(ctx context.Context) ([]*StudyProgram, error)
	GetStudyProgram(ctx context.Context, id primitive.ObjectID) (*StudyProgram, error)
	UpdateStudyProgram(ctx context.Context, id primitive.ObjectID, data *StudyProgram) error
	DeleteStudyProgram(ctx context.Context, id primitive.ObjectID) error
}

type studyPrgramRepository struct {
	collection *mongo.Collection
}

func NewStudyProgramRepository(collection *mongo.Collection) StudyProgramRepository {
	return &studyPrgramRepository{
		collection: collection,
	}
}

func (r *studyPrgramRepository) CreateStudyProgram(ctx context.Context, data *StudyProgram) error {

	filter := bson.M{
		"month":      data.Month,
		"year":       data.Year,
		"student_id": data.StudentID,
		"parent_id":  data.ParentID,
	}

	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	_, err := r.collection.InsertOne(ctx, data)
	return err

}

func (r *studyPrgramRepository) GetStudyPrograms(ctx context.Context) ([]*StudyProgram, error) {

	var result []*StudyProgram

	filter := bson.M{
		"is_deleted": false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *studyPrgramRepository) GetStudyProgram(ctx context.Context, id primitive.ObjectID) (*StudyProgram, error) {

	var result *StudyProgram

	filter := bson.M{
		"_id": id,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return result, nil

}

func (r *studyPrgramRepository) UpdateStudyProgram(ctx context.Context, id primitive.ObjectID, data *StudyProgram) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (r *studyPrgramRepository) DeleteStudyProgram(ctx context.Context, id primitive.ObjectID) error {
	
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
