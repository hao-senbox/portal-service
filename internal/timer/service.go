package timer

import (
	"context"
	"fmt"
	"portal/internal/user"
	"portal/pkg/uploader"
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
	ImageService    uploader.ImageService
}

func NewTimerService(TimerRepository TimerRepository, userService user.UserService, imageService uploader.ImageService) TimerService {
	return &timerService{
		TimerRepository: TimerRepository,
		UserService:     userService,
		ImageService:    imageService,
	}
}

func (s *timerService) CreateTimer(ctx context.Context, req *CreateTimerRequest, userID string) (string, error) {

	if req.Duration == 0 {
		return "", fmt.Errorf("duration is required")
	}

	timer := &Timer{
		ID:                primitive.NewObjectID(),
		StudentID:         req.StudentID,
		StartColor:        req.StartColor,
		EndColor:          req.EndColor,
		Duration:          req.Duration,
		NumberOfSound:     req.NumberOfSound,
		CenterLine:        req.LineCenter,
		OpacityDuration:   req.OpacityDuration,
		ImageStartKey:     req.ImageStartKey,
		ShowImageStart:    req.ShowImageStart,
		CaptionImageStart: req.CaptionImageStart,
		ImageEndKey:       req.ImageEndKey,
		ShowImageEnd:      req.ShowImageEnd,
		CaptionImageEnd:   req.CaptionImageEnd,
		TypePlay:          req.TypePlay,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
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

		if timer.ImageStartKey != "" {
			imageUrlStart, err := s.ImageService.GetImageKey(ctx, timer.ImageStartKey)
			if err != nil {
				return nil, err
			}
			timer.ImageStartKey = imageUrlStart.Url
		}

		if timer.ImageEndKey != "" {
			imageUrlEnd, err := s.ImageService.GetImageKey(ctx, timer.ImageEndKey)
			if err != nil {
				return nil, err
			}
			timer.ImageEndKey = imageUrlEnd.Url
		}

		result = append(result, &TimerResponse{
			ID:                timer.ID,
			Student:           student,
			StartColor:        timer.StartColor,
			EndColor:          timer.EndColor,
			Duration:          timer.Duration,
			NumberOfSound:     timer.NumberOfSound,
			CenterLine:        timer.CenterLine,
			OpacityDuration:   timer.OpacityDuration,
			ImageStartKey:     timer.ImageStartKey,
			ShowImageStart:    timer.ShowImageStart,
			CaptionImageStart: timer.CaptionImageStart,
			ImageEndKey:       timer.ImageEndKey,
			ShowImageEnd:      timer.ShowImageEnd,
			CaptionImageEnd:   timer.CaptionImageEnd,
			TypePlay:          timer.TypePlay,
			Teacher:           teacher,
			CreatedAt:         timer.CreatedAt,
			UpdatedAt:         timer.UpdatedAt,
		})
	}

	return result, nil
}
