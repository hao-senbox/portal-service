package program_planner

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramPlanerService interface {
	CreateProgramPlaner(ctx context.Context, req *CreateProgramPlanerRequest, userID string) (string, error)
	GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error)
	GetProgramPlaner(ctx context.Context, id string) (*ProgramPlaner, error)
	UpdateProgramPlaner(ctx context.Context, req *UpdateProgramPlanerRequest, id string) error
	DeleteProgramPlaner(ctx context.Context, id string) error
}

type programPlannerService struct {
	ProgramPlanerRepository ProgramPlanerRepository
}

func NewProgramPlanerService(ProgramPlanerRepository ProgramPlanerRepository) ProgramPlanerService {
	return &programPlannerService{
		ProgramPlanerRepository: ProgramPlanerRepository,
	}
}

func (service *programPlannerService) CreateProgramPlaner(ctx context.Context, req *CreateProgramPlanerRequest, userID string) (string, error) {

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}

	if req.Month == 0 && req.Year == 0 {
		return "", fmt.Errorf("month and year is required")
	}

	if len(req.SelectedSlots) == 0 {
		return "", fmt.Errorf("selected_slots is required")
	}

	var totalFee float64
	for _, slot := range req.SelectedSlots {
		totalFee += slot.Fee
	}

	programPlaner := &ProgramPlaner{
		ID:             primitive.NewObjectID(),
		StudentID:      req.StudentID,
		OrganizationID: req.OrganizationID,
		Month:          req.Month,
		Year:           req.Year,
		TotalFee:       totalFee,
		SelectedSlots:  req.SelectedSlots,
		Weeks:          []WeekPlan{},
		CreatedBy:      userID,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return service.ProgramPlanerRepository.CreateProgramPlaner(ctx, programPlaner)
}

func (service *programPlannerService) GetAllProgramPlaner(ctx context.Context) ([]*ProgramPlaner, error) {
	return service.ProgramPlanerRepository.GetAllProgramPlaner(ctx)
}

func (service *programPlannerService) GetProgramPlaner(ctx context.Context, id string) (*ProgramPlaner, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return service.ProgramPlanerRepository.GetProgramPlaner(ctx, objectID)
}

func (service *programPlannerService) UpdateProgramPlaner(ctx context.Context, req *UpdateProgramPlanerRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	program, err := service.ProgramPlanerRepository.GetProgramPlaner(ctx, objectID)
	if err != nil {
		return err
	}

	if req.Month != 0 {
		program.Month = req.Month
	}

	if req.Year != 0 {
		program.Year = req.Year
	}

	if len(req.SelectedSlots) != 0 {
		program.SelectedSlots = req.SelectedSlots
	}

	var totalFeeUpdate float64
	for _, slot := range program.SelectedSlots {
		totalFeeUpdate += slot.Fee
	}

	program.TotalFee = totalFeeUpdate

	return service.ProgramPlanerRepository.UpdateProgramPlaner(ctx, program, objectID)

}

func (service *programPlannerService) DeleteProgramPlaner(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return service.ProgramPlanerRepository.DeleteProgramPlaner(ctx, objectID)
}