package bmi

import (
	"context"
	"fmt"
	"math"
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BMIService interface {
	CreateBMI(ctx context.Context, req *CreateBMIStudentRequest, userID string) (string, error)
	GetBMIs(ctx context.Context, student_id string, date string) ([]*BMIStudentResponse, error)
	GetBMI(ctx context.Context, id string) (*BMIStudentResponse, error)
}

type bmiService struct {
	BMIRepo BMIRepo
	UserService user.UserService
}

func NewBMIService(BMIRepo BMIRepo, userService user.UserService) BMIService {
	return &bmiService{
		BMIRepo: BMIRepo,
		UserService: userService,
	}
}

func (s *bmiService) CreateBMI(ctx context.Context, req *CreateBMIStudentRequest, userID string) (string, error) {

	if req.Date == "" {
		return "", fmt.Errorf("date is required")
	}

	parseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	if req.Height == 0 {
		return "", fmt.Errorf("height is required")
	}

	if req.Weight == 0 {
		return "", fmt.Errorf("weight is required")
	}

	bmi := &BMI{
		ID:        primitive.NewObjectID(),
		StudentID: req.StudentID,
		Date:      parseDate,
		Height:    req.Height,
		Weight:    req.Weight,
		CreatedBy: userID,
		BMI:       calculateBMI(req.Height, req.Weight),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.BMIRepo.CreateBMI(ctx, bmi)
}

func (s *bmiService) GetBMIs(ctx context.Context, student_id string, date string) ([]*BMIStudentResponse, error) {
	
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

	bmis, err := s.BMIRepo.GetBMIs(ctx, student_id, dateRepo)

	if err != nil {
		return nil, err
	}

	var result []*BMIStudentResponse
	for _, bmi := range bmis {
		teacher, err := s.UserService.GetTeacherInfor(ctx, bmi.CreatedBy)
		if err != nil {
			return nil, err
		}

		student, err := s.UserService.GetStudentInfor(ctx, bmi.StudentID)
		if err != nil {
			return nil, err
		}

		result = append(result, &BMIStudentResponse{
			ID:        bmi.ID,
			Student:   student,
			Date:      bmi.Date.Format("2006-01-02"),
			BMI:       math.Round(bmi.BMI*100) / 100,
			Height:    bmi.Height,
			Weight:    bmi.Weight,
			Teacher:   teacher,
			CreatedAt: bmi.CreatedAt,
			UpdatedAt: bmi.UpdatedAt,
		})
	}

	return result, nil
}

func (s *bmiService) GetBMI(ctx context.Context, id string) (*BMIStudentResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	bmi, err := s.BMIRepo.GetBMI(ctx, objectID)
	if err != nil {
		return nil, err
	}

	teacher, err := s.UserService.GetTeacherInfor(ctx, bmi.CreatedBy)
	if err != nil {
		return nil, err
	}

	student, err := s.UserService.GetStudentInfor(ctx, bmi.StudentID)
	if err != nil {
		return nil, err
	}

	return &BMIStudentResponse{
		ID:        bmi.ID,
		Student:   student,
		Date:      bmi.Date.Format("2006-01-02"),
		BMI:       math.Round(bmi.BMI*100) / 100,
		Height:    bmi.Height,
		Weight:    bmi.Weight,
		Teacher:   teacher,
		CreatedAt: bmi.CreatedAt,
		UpdatedAt: bmi.UpdatedAt,
	}, nil
	
}

func calculateBMI(height float64, weight float64) float64 {
	height = height / 100
	return weight / (height * height)
}
