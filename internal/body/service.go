package body

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BodyService interface {
	CreateCheckIn(ctx context.Context, req *CreateCheckInRequest, userID string) (string, error)
}

type bodyService struct {
	BodyRepository BodyRepository
}

func NewBodyService(bodyRepository BodyRepository) BodyService {
	return &bodyService{
		BodyRepository: bodyRepository,
	}
}

func (s *bodyService) CreateCheckIn(ctx context.Context, req *CreateCheckInRequest, userID string) (string, error) {

	if userID == "" {
		return "", nil
	}

	if req.Date == "" {
		return "", fmt.Errorf("date is required")
	}

	parseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	if req.Context == "" {
		return "", fmt.Errorf("context is required")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.Marks) == 0 {
		return "", fmt.Errorf("marks is required")
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

	id, err := s.BodyRepository.CreateCheckIn(ctx, checkIn)

	if err != nil {
		return "", err
	}

	return id, nil
}
