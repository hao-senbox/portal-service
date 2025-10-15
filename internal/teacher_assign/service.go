package teacherassign

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherAssignmentService interface {
	CreateTeacherAssignment(ctx context.Context, req *CreateTeacherAssignmentRequest, userID string) (string, error)
	GetAllTeacherAssignment(ctx context.Context) ([]*TeacherAssignment, error)
	GetTeacherAssignment(ctx context.Context, id string) (*TeacherAssignment, error)
	UpdateTeacherAssignment(ctx context.Context, id string, req *UpdateTeacherAssignmentRequest) error
	DeleteTeacherAssignment(ctx context.Context, id string) error
}

type teacherAssignmentService struct {
	repository TeacherAssignmentRepository
}

func NewTeacherAssignmentService(repository TeacherAssignmentRepository) TeacherAssignmentService {
	return &teacherAssignmentService{
		repository: repository,
	}
}

func (s *teacherAssignmentService) CreateTeacherAssignment(ctx context.Context, req *CreateTeacherAssignmentRequest, userID string) (string, error) {

	if userID == "" {
		return "", errors.New("user ID is empty")
	}

	if req.ParentID == "" {
		return "", errors.New("parent ID is empty")
	}

	if req.StudentID == "" {
		return "", errors.New("student ID is empty")
	}

	var total float64

	if req.AgeRange == nil {
		req.AgeRange = &Data{}
	}

	if req.SkillSet == nil {
		req.SkillSet = &Data{}
	}

	if req.PDLevel == nil {
		req.PDLevel = &Data{}
	}

	if req.ExperienceExt == nil {
		req.ExperienceExt = &Data{}
	}

	if req.ExperienceInt == nil {
		req.ExperienceInt = &Data{}
	}

	if req.Qualification == nil {
		req.Qualification = &Data{}
	}

	if req.Language == nil {
		req.Language = &Data{}
	}

	items := []*Data{
		req.AgeRange,
		req.SkillSet,
		req.PDLevel,
		req.ExperienceExt,
		req.ExperienceInt,
		req.Qualification,
		req.Language,
	}

	for _, item := range items {
		total += item.Price
	}

	ID := primitive.NewObjectID()

	data := &TeacherAssignment{
		ID:            ID,
		ParentID:      req.ParentID,
		StudentID:     req.StudentID,
		Month:         req.Month,
		Year:          req.Year,
		Language:      req.Language,
		Qualification: req.Qualification,
		ExperienceExt: req.ExperienceExt,
		ExperienceInt: req.ExperienceInt,
		PDLevel:       req.PDLevel,
		SkillSet:      req.SkillSet,
		AgeRange:      req.AgeRange,
		MonthlyFee:    total,
		CreatedBy:     userID,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return ID.Hex(), s.repository.CreateTeacherAssignment(ctx, data)
}

func (s *teacherAssignmentService) GetAllTeacherAssignment(ctx context.Context) ([]*TeacherAssignment, error) {
	return s.repository.GetAllTeacherAssignment(ctx)
}

func (s *teacherAssignmentService) GetTeacherAssignment(ctx context.Context, id string) (*TeacherAssignment, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repository.GetTeacherAssignment(ctx, objectId)
}

func (s *teacherAssignmentService) UpdateTeacherAssignment(ctx context.Context, id string, req *UpdateTeacherAssignmentRequest) error {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	teacherAssign, err := s.repository.GetTeacherAssignment(ctx, objectId)
	if err != nil {
		return err
	}

	if req.Language != nil {
		teacherAssign.Language = req.Language
	}

	if req.Qualification != nil {
		teacherAssign.Qualification = req.Qualification
	}

	if req.ExperienceExt != nil {
		teacherAssign.ExperienceExt = req.ExperienceExt
	}

	if req.ExperienceInt != nil {
		teacherAssign.ExperienceInt = req.ExperienceInt
	}

	if req.PDLevel != nil {
		teacherAssign.PDLevel = req.PDLevel
	}

	if req.SkillSet != nil {
		teacherAssign.SkillSet = req.SkillSet
	}

	if req.AgeRange != nil {
		teacherAssign.AgeRange = req.AgeRange
	}

	items := []*Data{
		teacherAssign.Language,
		teacherAssign.Qualification,
		teacherAssign.ExperienceExt,
		teacherAssign.ExperienceInt,
		teacherAssign.PDLevel,
		teacherAssign.SkillSet,
		teacherAssign.AgeRange,
	}

	var total float64

	for _, item := range items {
		total += item.Price
	}

	teacherAssign.MonthlyFee = total

	teacherAssign.UpdatedAt = time.Now()
	return s.repository.UpdateTeacherAssignment(ctx, objectId, teacherAssign)
}

func (s *teacherAssignmentService) DeleteTeacherAssignment(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repository.DeleteTeacherAssignment(ctx, objectId)
}
