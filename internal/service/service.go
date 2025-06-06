package service

import (
	"context"
	"fmt"
	"portal/internal/models"
	"portal/internal/repository"
	"time"
)

type PortalService interface {
	CreateStudentActivity(ctx context.Context, req *models.RequestStudentActivity) error
	GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*models.StudentActivity, error)
}
type portalService struct {
	repoPortal repository.PortalRepository
}

func NewPortalService(repo repository.PortalRepository) PortalService {
	return &portalService{
		repoPortal: repo,
	}
}

func (s *portalService) CreateStudentActivity(ctx context.Context, req *models.RequestStudentActivity) error {

	if req.StudentID == "" {
		return fmt.Errorf("student ID cannot be empty")
	}	

	if req.TypeActivity == "" {
		return fmt.Errorf("type activity cannot be empty")
	}

	if req.Date == "" {
		return fmt.Errorf("date cannot be empty")
	}

	dateParse, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	if req.Data == nil {	
		return fmt.Errorf("data cannot be empty")
	}

	if req.AssignedBy == "" {
		return fmt.Errorf("assigned by cannot be empty")
	}

	studentActivity := &models.StudentActivity{
		StudentID:    req.StudentID,
		TypeActivity: req.TypeActivity,
		Date:         dateParse,
		Data:         req.Data,
		SubittedAt:   time.Now(),
		CreatdBy:     req.AssignedBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),	
	}

	err = s.repoPortal.CreateStudentActivity(ctx, studentActivity)
	if err != nil {
		return fmt.Errorf("failed to create student activity: %w", err)
	}

	return nil
}

func (s *portalService) GetAllStudentActivity(ctx context.Context, studentID string, date string) ([]*models.StudentActivity, error) {

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	
	activities, err := s.repoPortal.GetAllStudentActivity(ctx, studentID, parsedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activity: %w", err)
	}

	return activities, nil
}
