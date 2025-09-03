package body

import (
	"context"
	"fmt"
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BodyService interface {
	CreateCheckIn(ctx context.Context, req *CreateCheckInRequest, userID string) error
	GetCheckIns(ctx context.Context, student_id string, date string) ([]*CheckInReponse, error)
}

type bodyService struct {
	BodyRepository BodyRepository
	UserService    user.UserService
}

func NewBodyService(bodyRepository BodyRepository, userService user.UserService) BodyService {
	return &bodyService{
		BodyRepository: bodyRepository,
		UserService:    userService,
	}
}

func (s *bodyService) CreateCheckIn(ctx context.Context, req *CreateCheckInRequest, userID string) error {

	if userID == "" {
		return nil
	}

	if req.Date == "" {
		return fmt.Errorf("date is required")
	}

	parseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	if req.Context == "" {
		return fmt.Errorf("context is required")
	}

	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	if req.StudentID == "" {
		return fmt.Errorf("student_id is required")
	}

	if len(req.Marks) == 0 {
		return fmt.Errorf("marks is required")
	}

	now := time.Now()
	for i := range req.Marks {
		req.Marks[i].SubmittedAt = now
	}

	checkIn := &CheckIn{
		ID:        primitive.NewObjectID(),
		StudentID: req.StudentID,
		Date:      parseDate,
		Context:   req.Context,
		Gender:    req.Gender,
		Type:      req.Type,
		Marks:     req.Marks,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.BodyRepository.PushCheckIn(ctx, checkIn)

	if err != nil {
		return err
	}

	return nil
}

func (s *bodyService) GetCheckIns(ctx context.Context, student_id string, date string) ([]*CheckInReponse, error) {

	var dateRepo *time.Time

	if date != "" {
		parseDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		dateRepo = &parseDate
	} else {
		dateRepo = nil
	}

	checkIns, err := s.BodyRepository.GetCheckIns(ctx, student_id, dateRepo)
	if err != nil {
		return nil, err
	}
	
	var result []*CheckInReponse

	for _, checkIn := range checkIns {
		student, err := s.UserService.GetStudentInfor(ctx, checkIn.StudentID)
		if err != nil {
			return nil, err
		}

		teacher, err := s.UserService.GetTeacherInfor(ctx, checkIn.CreatedBy)
		if err != nil {
			return nil, err
		}

		result = append(result, &CheckInReponse{
			ID:        checkIn.ID,
			Student:   student,
			Date:      checkIn.Date.Format("2006-01-02"),
			Context:   checkIn.Context,
			Type:      checkIn.Type,
			Gender:    checkIn.Gender,
			Marks:     checkIn.Marks,
			Teacher:   teacher,
			CreatedAt: checkIn.CreatedAt,
			UpdatedAt: checkIn.UpdatedAt,
		})
	}

	return result, nil

}
