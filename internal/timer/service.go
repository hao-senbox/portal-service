package timer

import (
	"context"
	"fmt"
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimerService interface {
	CreateTimer(ctx context.Context, req *CreateTimerRequest, userID string) (string, error)
	GetTimers(ctx context.Context, studentID string) ([]*TimerResponse, error)
}

type timerService struct {
	TimerRepository TimerRepository
	UserService     user.UserService
}

func NewTimerService(TimerRepository TimerRepository, userService user.UserService) TimerService {
	return &timerService{
		TimerRepository: TimerRepository,
		UserService:     userService,
	}
}

func (s *timerService) CreateTimer(ctx context.Context, req *CreateTimerRequest, userID string) (string, error) {

	if req.Duration == 0 {
		return "", fmt.Errorf("duration is required")
	}

	timer := &Timer{
		ID:            primitive.NewObjectID(),
		StudentID:     req.StudentID,
		StartColor:    req.StartColor,
		EndColor:      req.EndColor,
		Duration:      req.Duration,
		NumberOfSound: req.NumberOfSound,
		Image:         req.Image,
		TypePlay:      req.TypePlay,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	id, err := s.TimerRepository.CreateTimer(ctx, timer)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (s *timerService) GetTimers(ctx context.Context, studentID string) ([]*TimerResponse, error) {

	if studentID == "" {
		studentID = ""
	}

	timers, err := s.TimerRepository.GetTimers(ctx, studentID)

	if err != nil {
		return nil, err
	}

	var result []*TimerResponse
	for _, timer := range timers {

		student, err := s.UserService.GetStudentInfor(ctx, timer.StudentID)
		if err != nil {
			return nil, err
		}

		teacher, err := s.UserService.GetUserInfor(ctx, timer.CreatedBy)
		if err != nil {
			return nil, err
		}

		result = append(result, &TimerResponse{
			ID:            timer.ID,
			Student:       student,
			StartColor:    timer.StartColor,
			EndColor:      timer.EndColor,
			Duration:      timer.Duration,
			NumberOfSound: timer.NumberOfSound,
			Image:         timer.Image,
			TypePlay:      timer.TypePlay,
			Teacher:       teacher,
			CreatedAt:     timer.CreatedAt,
			UpdatedAt:     timer.UpdatedAt,
		})
	}

	return result, nil
}
