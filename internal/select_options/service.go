package selectoptions

import (
	"context"
	"fmt"
	"portal/internal/term"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectOptionsService interface {
	CreateSelectOption(ctx context.Context, req CreateSelectOptionRequest, userID string) (string, error)
}

type selectOptionsService struct {
	repo        SelectOptionsRepository
	termService term.TermService
}

func NewSelectOptionsService(repo SelectOptionsRepository, termService term.TermService) SelectOptionsService {
	return &selectOptionsService{
		repo:        repo,
		termService: termService,
	}
}

func (s *selectOptionsService) CreateSelectOption(ctx context.Context, req CreateSelectOptionRequest, userID string) (string, error) {

	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.Options) == 0 {
		return "", fmt.Errorf("options array cannot be empty")
	}

	var options []Options

	if req.Type == "select_topic" {
		options = []Options{
			{TopicID: "6909c4192c3eb1f1800d757d", Order: 1},
			{TopicID: "6909c4d5c3189f6c406ba20c", Order: 2},
			{TopicID: "6909c4f9c3189f6c406ba20d", Order: 3},
			{TopicID: "6909c515c3189f6c406ba20e", Order: 4},
		}
		if len(req.Options) > 0 {
			options = append(options, req.Options...)
		}
	} else {
		options = req.Options
	}

	for i, opt := range req.Options {
		if opt.Order == 0 {
			opt.Order = i + 1
			req.Options[i] = opt
		}
	}

	term, err := s.termService.GetCurrentTermByOrgID(ctx, req.OrganizationID)
	if err != nil {
		return "", fmt.Errorf("failed to get current term by org id: %w", err)
	}

	if term == nil {
		return "", fmt.Errorf("current term not found")
	}

	id := primitive.NewObjectID()

	doc := &SelectOptions{
		ID:             id,
		OrganizationID: req.OrganizationID,
		StudentID:      req.StudentID,
		TermID:         term.ID,
		Type:           req.Type,
		Options:        options,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsDeleted:      false,
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return "", fmt.Errorf("failed to create select option: %w", err)
	}

	return id.Hex(), nil
}
