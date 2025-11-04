package studypreference

import (
	"context"
	"fmt"
	"portal/internal/term"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyPreferenceService interface {
	CreateStudyPreference(ctx context.Context, req *CreateStudyPreferenceRequest, userID string) (string, error)
	GetStudyPreferencesByStudentID(ctx context.Context, studentID, orgID string) (*StudyPreference, error)
	GetStudyPreferenceByID(ctx context.Context, id string) (*StudyPreference, error)
	UpdateStudyPreference(ctx context.Context, id string, req *UpdateStudyPreferenceRequest, userID string) error
}

type studyPreferenceService struct {
	studyPreferenceRepository StudyPreferenceRepository
	termService term.TermService
}

func NewStudyPreferenceService(studyPreferenceRepository StudyPreferenceRepository, termService term.TermService) StudyPreferenceService {
	return &studyPreferenceService{studyPreferenceRepository: studyPreferenceRepository, termService: termService}
}

func (s *studyPreferenceService) CreateStudyPreference(ctx context.Context, req *CreateStudyPreferenceRequest, userID string) (string, error) {

	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}

	if req.TermID == "" {
		return "", fmt.Errorf("term_id is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.TeacherSelection) == 0 {
		return "", fmt.Errorf("teacher_selection is required")
	}

	data := &StudyPreference{
		ID:               primitive.NewObjectID(),
		OrganizationID:   req.OrganizationID,
		TermID:           req.TermID,
		StudentID:        req.StudentID,
		ParentSelection:  []Data{},
		TeacherSelection: req.TeacherSelection,
		CreatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return data.ID.Hex(), s.studyPreferenceRepository.CreateStudyPreference(ctx, data)

}

func (s *studyPreferenceService) GetStudyPreferencesByStudentID(ctx context.Context, studentID, orgID string) (*StudyPreference, error) {
	if studentID == "" {
		return nil, fmt.Errorf("student_id is required")
	}

	if orgID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	term, err := s.termService.GetCurrentTermByOrgID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current term by org id: %w", err)
	}

	if term == nil {
		return nil, fmt.Errorf("current term not found")
	}

	data, err := s.studyPreferenceRepository.GetStudyPreferencesByStudentID(ctx, studentID, term.ID, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study preferences by student id: %w", err)
	}

	for i := range data.TeacherSelection {
		if len(data.TeacherSelection[i].Pairs) == 2 {
			pairs := data.TeacherSelection[i].Pairs
			if pairs[0].Value > pairs[1].Value {
				data.TeacherSelection[i].Selected = pairs[0].Category
			} else if pairs[1].Value > pairs[0].Value {
				data.TeacherSelection[i].Selected = pairs[1].Category
			} else {
				data.TeacherSelection[i].Selected = "equal"
			}
		}
	}

	for i := range data.ParentSelection {
		if len(data.ParentSelection[i].Pairs) == 2 {
			pairs := data.ParentSelection[i].Pairs
			if pairs[0].Value > pairs[1].Value {
				data.ParentSelection[i].Selected = pairs[0].Category
			} else if pairs[1].Value > pairs[0].Value {
				data.ParentSelection[i].Selected = pairs[1].Category
			} else {
				data.ParentSelection[i].Selected = "equal"
			}
		}
	}

	return data, nil
}


func (s *studyPreferenceService) GetStudyPreferenceByID(ctx context.Context, id string) (*StudyPreference, error) {

	if id == "" {
		return nil, fmt.Errorf("study preference id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid study preference id format: %w", err)
	}

	data, err := s.studyPreferenceRepository.GetStudyPreferenceByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get study preference by id: %w", err)
	}

	for i := range data.TeacherSelection {
		if len(data.TeacherSelection[i].Pairs) == 2 {
			pairs := data.TeacherSelection[i].Pairs
			if pairs[0].Value > pairs[1].Value {
				data.TeacherSelection[i].Selected = pairs[0].Category
			} else if pairs[1].Value > pairs[0].Value {
				data.TeacherSelection[i].Selected = pairs[1].Category
			} else {
				data.TeacherSelection[i].Selected = "equal"
			}
		}
	}

	for i := range data.ParentSelection {
		if len(data.ParentSelection[i].Pairs) == 2 {
			pairs := data.ParentSelection[i].Pairs
			if pairs[0].Value > pairs[1].Value {
				data.ParentSelection[i].Selected = pairs[0].Category
			} else if pairs[1].Value > pairs[0].Value {
				data.ParentSelection[i].Selected = pairs[1].Category
			} else {
				data.ParentSelection[i].Selected = "equal"
			}
		}
	}

	return data, nil
}

func (s *studyPreferenceService) UpdateStudyPreference(ctx context.Context, id string, req *UpdateStudyPreferenceRequest, userID string) error {

	if id == "" {
		return fmt.Errorf("study preference id is required")
	}

	if len(req.ParentSelections) == 0 {
		return fmt.Errorf("parent_selections cannot be empty")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid study preference id format: %w", err)
	}

	updateData := map[string]interface{}{
		"parent_selections": req.ParentSelections,
		"updated_by":        userID,
		"updated_at":        time.Now(),
	}

	return s.studyPreferenceRepository.UpdateStudyPreference(ctx, objectID, updateData)
}
