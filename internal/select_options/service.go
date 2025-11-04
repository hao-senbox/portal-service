package selectoptions

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectOptionsService interface {
	CreateSelectOption(ctx context.Context, req CreateSelectOptionRequest, userID string) (string, error)
}

type selectOptionsService struct {
	repo SelectOptionsRepository
}

func NewSelectOptionsService(repo SelectOptionsRepository) SelectOptionsService {
	return &selectOptionsService{
		repo: repo,
	}
}

func (s *selectOptionsService) CreateSelectOption(ctx context.Context, req CreateSelectOptionRequest, userID string) (string, error) {

	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	if req.TermID == "" {
		return "", fmt.Errorf("term_id is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.Options) == 0 {
		return "", fmt.Errorf("options array cannot be empty")
	}

	for i, opt := range req.Options {
		if opt.Name == "" {
			return "", fmt.Errorf("options[%d].name is required", i)
		}
		if opt.Order == 0 {
			opt.Order = i + 1
			req.Options[i] = opt
		}
	}

	id := primitive.NewObjectID()

	doc := &SelectOptions{
		ID:             id,
		OrganizationID: req.OrganizationID,
		StudentID:      req.StudentID,
		TermID:         req.TermID,
		Type:           req.Type,
		Options:        req.Options,
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
