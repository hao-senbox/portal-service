package teacherassign

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeacherAssignmentRepository interface {
	CreateTeacherAssignment(ctx context.Context, data *TeacherAssignment) error
	GetAllTeacherAssignment(ctx context.Context) ([]*TeacherAssignment, error)
	GetTeacherAssignment(ctx context.Context, id primitive.ObjectID) (*TeacherAssignment, error)
	UpdateTeacherAssignment(ctx context.Context, id primitive.ObjectID, data *TeacherAssignment) error
	DeleteTeacherAssignment(ctx context.Context, id primitive.ObjectID) error
}

type teacherAssignmentRepository struct {
	collection *mongo.Collection
}

func NewTeacherAssignmentRepository(collection *mongo.Collection) TeacherAssignmentRepository {
	return &teacherAssignmentRepository{
		collection: collection,
	}
}

func (r *teacherAssignmentRepository) CreateTeacherAssignment(ctx context.Context, data *TeacherAssignment) error {

	filter := bson.M{
		"month": data.Month,
		"year":  data.Year,
		"parent_id": data.ParentID,
		"student_id": data.StudentID,
	}

	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	_, err := r.collection.InsertOne(ctx, data)

	return err
}

func (r *teacherAssignmentRepository) GetAllTeacherAssignment(ctx context.Context) ([]*TeacherAssignment, error) {

	var teacherAssignments []*TeacherAssignment

	filter := bson.M{
		"is_deleted": false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &teacherAssignments); err != nil {
		return nil, err
	}

	return teacherAssignments, nil

}

func (r *teacherAssignmentRepository) GetTeacherAssignment(ctx context.Context, id primitive.ObjectID) (*TeacherAssignment, error) {
	var teacherAssignment TeacherAssignment

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&teacherAssignment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &teacherAssignment, nil

}

func (r *teacherAssignmentRepository) UpdateTeacherAssignment(ctx context.Context, id primitive.ObjectID, data *TeacherAssignment) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (r *teacherAssignmentRepository) DeleteTeacherAssignment(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_deleted": true}})
	return err
}