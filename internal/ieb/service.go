package ieb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEBService interface {
	CreateIEB(ctx context.Context, req *CreateIEBRequest, userID string) (string, error)
	GetIEB(ctx context.Context, userID string, termID string, languageKey string, regionKey string) (*IEB, error)
}

type iebService struct {
	iebRepository IEBRepository
}

func NewIEBService(iebRepository IEBRepository) IEBService {
	return &iebService{
		iebRepository: iebRepository,
	}
}

func (service *iebService) CreateIEB(ctx context.Context, req *CreateIEBRequest, userID string) (string, error) {

	if req.TermID == "" {
		return "", fmt.Errorf("term_id is required")
	}

	if req.Owner.OwnerID == "" {
		return "", fmt.Errorf("owner_id is required")
	}

	if req.LanguageKey == "" {
		return "", fmt.Errorf("language_id is required")
	}

	if req.RegionKey == "" {
		return "", fmt.Errorf("region_id is required")
	}

	ID := primitive.NewObjectID()

	data := &IEB{
		ID:          ID,
		Owner:       req.Owner,
		TermID:      req.TermID,
		LanguageKey: req.LanguageKey,
		RegionKey:   req.RegionKey,
		Information: req.Information,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return ID.Hex(), service.iebRepository.CreateIEB(ctx, data)

}

func (service *iebService) GetIEB(ctx context.Context, userID string, termID string, languageKey string, regionKey string) (*IEB, error) {

	if termID == "" {
		return nil, fmt.Errorf("term_id is required")
	}

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if languageKey == "" {
		return nil, fmt.Errorf("language_id is required")
	}

	if regionKey == "" {
		return nil, fmt.Errorf("region_id is required")
	}

	ieb, err := service.iebRepository.GetIEB(ctx, userID, termID, languageKey, regionKey)
	if err != nil {
		return nil, nil
	}

	return ieb, nil
	
}
