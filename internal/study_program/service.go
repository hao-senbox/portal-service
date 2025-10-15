package studyprogram

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyProgramService interface {
	CreateStudyProgram(ctx context.Context, req *CreateStudyProgramRequest, userID string) (string, error)
	GetStudyPrograms(ctx context.Context) ([]*StudyProgram, error)
	GetStudyProgram(ctx context.Context, id string) (*StudyProgram, error)
	UpdateStudyProgram(ctx context.Context, id string, req *UpdateStudyProgramRequest) error
	DeleteStudyProgram(ctx context.Context, id string) error
}

type studyProgramService struct {
	StudyProgramRepository StudyProgramRepository
}

func NewStudyProgramService(repository StudyProgramRepository) StudyProgramService {
	return &studyProgramService{
		StudyProgramRepository: repository,
	}
}

func (s *studyProgramService) CreateStudyProgram(ctx context.Context, req *CreateStudyProgramRequest, userID string) (string, error) {

	if userID == "" {
		return "", errors.New("user ID is empty")
	}

	if req.ParentID == "" {
		return "", errors.New("parent ID is empty")
	}

	if req.StudentID == "" {
		return "", errors.New("student ID is empty")
	}

	if req.TimeSlot == nil {
		req.TimeSlot = &Data{}
	}

	if req.ServiceRatio == nil {
		req.ServiceRatio = &Data{}
	}

	if req.SkillPercent == nil {
		req.SkillPercent = &Data{}
	}

	if req.TeacherWeight == nil {
		req.TeacherWeight = &Data{}
	}

	total := 0.0

	items := []Data{
		*req.TimeSlot,
		*req.ServiceRatio,
		*req.SkillPercent,
		*req.TeacherWeight,
	}

	for _, item := range items {
		total += item.Price
	}

	if req.Extras == nil {
		req.Extras = []*Data{}
	}

	if req.Extras != nil {
		for _, extra := range req.Extras {
			total += extra.Price
		}
	}

	if req.OtherFees == nil {
		req.OtherFees = []*Data{}
	}

	if req.OtherFees != nil {
		for _, fee := range req.OtherFees {
			total += fee.Price
		}
	}

	ID := primitive.NewObjectID()

	data := &StudyProgram{
		ID:            ID,
		ParentID:      req.ParentID,
		StudentID:     req.StudentID,
		Month:         req.Month,
		Year:          req.Year,
		TimeSlot:      req.TimeSlot,
		ServiceRatio:  req.ServiceRatio,
		SkillPercent:  req.SkillPercent,
		TeacherWeight: req.TeacherWeight,
		Extras:        req.Extras,
		OtherFees:     req.OtherFees,
		MonthlyTotal:  total,
		CreatedBy:     userID,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return ID.Hex(), s.StudyProgramRepository.CreateStudyProgram(ctx, data)
}

func (s *studyProgramService) GetStudyPrograms(ctx context.Context) ([]*StudyProgram, error) {
	return s.StudyProgramRepository.GetStudyPrograms(ctx)
}

func (s *studyProgramService) GetStudyProgram(ctx context.Context, id string) (*StudyProgram, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.StudyProgramRepository.GetStudyProgram(ctx, objectID)
}

func (s *studyProgramService) UpdateStudyProgram(ctx context.Context, id string, req *UpdateStudyProgramRequest) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	studyProgram, err := s.StudyProgramRepository.GetStudyProgram(ctx, objectID)
	if err != nil {
		return err
	}

	if studyProgram == nil {
		return errors.New("study program not found")
	}

	if req.Extras != nil {
		for _, extra := range req.Extras {
			if extra.Label == "" {
				return errors.New("label cannot be empty")
			}
			if extra.Price == 0 {
				return errors.New("price cannot be 0")
			}
		}

		studyProgram.Extras = req.Extras
	}

	if req.OtherFees != nil {
		for _, fee := range req.OtherFees {
			if fee.Label == "" {
				return errors.New("label cannot be empty")
			}
			if fee.Price == 0 {
				return errors.New("price cannot be 0")
			}
		}

		studyProgram.OtherFees = req.OtherFees
	}

	if req.ServiceRatio != nil {
		if req.ServiceRatio.Label == "" {
			return errors.New("label cannot be empty")
		}
		if req.ServiceRatio.Price == 0 {
			return errors.New("price cannot be 0")
		}
		studyProgram.ServiceRatio = req.ServiceRatio
	}

	if req.SkillPercent != nil {
		if req.SkillPercent.Label == "" {
			return errors.New("label cannot be empty")
		}
		if req.SkillPercent.Price == 0 {
			return errors.New("price cannot be 0")
		}
		studyProgram.SkillPercent = req.SkillPercent
	}

	if req.TeacherWeight != nil {
		if req.TeacherWeight.Label == "" {
			return errors.New("label cannot be empty")
		}
		if req.TeacherWeight.Price == 0 {
			return errors.New("price cannot be 0")
		}
		studyProgram.TeacherWeight = req.TeacherWeight
	}

	if req.TimeSlot != nil {
		if req.TimeSlot.Label == "" {
			return errors.New("label cannot be empty")
		}
		if req.TimeSlot.Price == 0 {
			return errors.New("price cannot be 0")
		}
		studyProgram.TimeSlot = req.TimeSlot
	}

	total := 0.0

	items := []*Data{
		studyProgram.TimeSlot,
		studyProgram.ServiceRatio,
		studyProgram.SkillPercent,
		studyProgram.TeacherWeight,
	}
	for _, item := range items {
		if item.Price > 0 {
			total += item.Price
		}
	}

	for _, extra := range studyProgram.Extras {
		if extra.Price > 0 {
			total += extra.Price
		}
	}
	for _, fee := range studyProgram.OtherFees {
		if fee.Price > 0 {
			total += fee.Price
		}
	}

	studyProgram.MonthlyTotal = total
	studyProgram.UpdatedAt = time.Now()

	return s.StudyProgramRepository.UpdateStudyProgram(ctx, objectID, studyProgram)
}

func (s *studyProgramService) DeleteStudyProgram(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.StudyProgramRepository.DeleteStudyProgram(ctx, objectID)

}
