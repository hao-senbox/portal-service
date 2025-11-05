package studypreference

import (
	"context"
	"fmt"
	selectoptions "portal/internal/select_options"
	"portal/internal/term"
	"portal/internal/topic"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyPreferenceService interface {
	CreateStudyPreference(ctx context.Context, req *CreateStudyPreferenceRequest, userID string) (string, error)
	GetStudyPreferencesByStudentID(ctx context.Context, studentID, orgID string) (*StudyPreference, error)
	// GetStudyPreferenceByID(ctx context.Context, id string) (*StudyPreference, error)
	UpdateStudyPreference(ctx context.Context, id string, req *UpdateStudyPreferenceRequest, userID string) error
	GetStudyPreferenceStatistical(ctx context.Context, orgID, studentID string) (map[string]interface{}, error)
}

type studyPreferenceService struct {
	studyPreferenceRepository StudyPreferenceRepository
	termService               term.TermService
	selectOptionsRepository   selectoptions.SelectOptionsRepository
	topicService              topic.TopicService
}

func NewStudyPreferenceService(
	studyPreferenceRepository StudyPreferenceRepository,
	termService term.TermService,
	selectOptionsRepository selectoptions.SelectOptionsRepository,
	topicService topic.TopicService,
) StudyPreferenceService {
	return &studyPreferenceService{
		studyPreferenceRepository: studyPreferenceRepository,
		termService:               termService,
		selectOptionsRepository:   selectOptionsRepository,
		topicService:              topicService,
	}
}

func (s *studyPreferenceService) CreateStudyPreference(ctx context.Context, req *CreateStudyPreferenceRequest, userID string) (string, error) {

	if req.OrganizationID == "" {
		return "", fmt.Errorf("organization_id is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student_id is required")
	}

	if len(req.TeacherSelection) == 0 {
		return "", fmt.Errorf("teacher_selection is required")
	}

	term, err := s.termService.GetCurrentTermByOrgID(ctx, req.OrganizationID)
	if err != nil {
		return "", fmt.Errorf("failed to get current term by org id: %w", err)
	}

	if term == nil {
		return "", fmt.Errorf("current term not found")
	}

	data := &StudyPreference{
		ID:               primitive.NewObjectID(),
		OrganizationID:   req.OrganizationID,
		TermID:           term.ID,
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
		return nil, fmt.Errorf("student id is required")
	}

	if orgID == "" {
		return nil, fmt.Errorf("org id is required")
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
		return nil, fmt.Errorf("failed to get study preference by id: %w", err)
	}

	top := []string{
		"Managing Self",
		"Gross Motor",
		"Writing",
		"Self Help",
		"Contributing Arts",
		"Group Work",
	}

	bottom := []string{
		"Socialize",
		"Fine Motor",
		"Drawing",
		"Life Skill",
		"Participating",
		"Independence",
	}

	isInArray := func(arr []string, val string) bool {
		for _, v := range arr {
			if strings.EqualFold(v, val) {
				return true
			}
		}
		return false
	}

	// ✅ xử lý TeacherSelection
	for i := range data.TeacherSelection {
		if len(data.TeacherSelection[i].Pairs) == 2 {
			pairs := data.TeacherSelection[i].Pairs
			p1, p2 := pairs[0], pairs[1]

			var chosen Pair
			if p1.Value > p2.Value {
				chosen = p1
			} else if p2.Value > p1.Value {
				chosen = p2
			} else {
				data.TeacherSelection[i].Selector = map[string]interface{}{
					"name":  "equal",
					"value": 50,
				}
				continue
			}

			var selectorName string
			if isInArray(top, chosen.Category) {
				selectorName = "bottom"
			} else if isInArray(bottom, chosen.Category) {
				selectorName = "top"
			} else {
				selectorName = "unknown"
			}

			data.TeacherSelection[i].Selector = map[string]interface{}{
				"name":  selectorName,
				"value": 100 - chosen.Value,
			}
		}
	}

	for i := range data.ParentSelection {
		if len(data.ParentSelection[i].Pairs) == 2 {
			pairs := data.ParentSelection[i].Pairs
			p1, p2 := pairs[0], pairs[1]

			var chosen Pair
			if p1.Value > p2.Value {
				chosen = p1
			} else if p2.Value > p1.Value {
				chosen = p2
			} else {
				data.ParentSelection[i].Selector = map[string]interface{}{
					"name":  "equal",
					"value": 50,
				}
				continue
			}

			var selectorName string
			if isInArray(top, chosen.Category) {
				selectorName = "bottom"
			} else if isInArray(bottom, chosen.Category) {
				selectorName = "top"
			} else {
				selectorName = "unknown"
			}

			data.ParentSelection[i].Selector = map[string]interface{}{
				"name":  selectorName,
				"value": 100 - chosen.Value,
			}
		}
	}

	return data, nil
}



// func (s *studyPreferenceService) GetStudyPreferenceByID(ctx context.Context, id string) (*StudyPreference, error) {

// 	if id == "" {
// 		return nil, fmt.Errorf("study preference id is required")
// 	}

// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid study preference id format: %w", err)
// 	}

// 	data, err := s.studyPreferenceRepository.GetStudyPreferenceByID(ctx, objectID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get study preference by id: %w", err)
// 	}

// 	for i := range data.TeacherSelection {
// 		if len(data.TeacherSelection[i].Pairs) == 2 {
// 			pairs := data.TeacherSelection[i].Pairs
// 			if pairs[0].Value > pairs[1].Value {
// 				data.TeacherSelection[i] = pairs[0].Category
// 			} else if pairs[1].Value > pairs[0].Value {
// 				data.TeacherSelection[i] = pairs[1].Category
// 			} else {
// 				data.TeacherSelection[i] = "equal"
// 			}
// 		}
// 	}

// 	for i := range data.ParentSelection {
// 		if len(data.ParentSelection[i].Pairs) == 2 {
// 			pairs := data.ParentSelection[i].Pairs
// 			if pairs[0].Value > pairs[1].Value {
// 				data.ParentSelection[i] = pairs[0].Category
// 			} else if pairs[1].Value > pairs[0].Value {
// 				data.ParentSelection[i] = pairs[1].Category
// 			} else {
// 				data.ParentSelection[i] = "equal"
// 			}
// 		}
// 	}

// 	return data, nil
// }

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

func (s *studyPreferenceService) GetStudyPreferenceStatistical(ctx context.Context, orgID, studentID string) (map[string]interface{}, error) {

	term, err := s.termService.GetCurrentTermByOrgID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current term by org id: %w", err)
	}

	if term == nil {
		return nil, fmt.Errorf("current term not found")
	}

	iepPriority, err := s.selectOptionsRepository.GetSelectOption(ctx, "iep_priority", orgID, term.ID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get iep priority: %w", err)
	}

	topicPlanner, err := s.selectOptionsRepository.GetSelectOption(ctx, "topic_planner", orgID, term.ID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic planner: %w", err)
	}

	lifeSkills, err := s.selectOptionsRepository.GetSelectOption(ctx, "life_skills", orgID, term.ID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get life skills: %w", err)
	}

	selectTopic, err := s.selectOptionsRepository.GetSelectOption(ctx, "select_topic", orgID, term.ID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get select topic: %w", err)
	}

	result := make(map[string]interface{})

	if iepPriority != nil && len(iepPriority.Options) > 0 {
		var names []string
		for _, opt := range iepPriority.Options {
			if opt.Name != "" {
				names = append(names, opt.Name)
			}
		}
		if len(names) > 0 {
			result["iep_priority"] = "" + strings.Join(names, ", ")
		}
	} else {
		result["iep_priority"] = "No iep priority"
	}

	if topicPlanner != nil && len(topicPlanner.Options) > 0 {
		var topics []string
		for _, opt := range topicPlanner.Options {
			if opt.Name != "" && opt.Status != nil {
				topics = append(topics, fmt.Sprintf("%s → %s", opt.Name, *opt.Status))
			}
		}
		if len(topics) > 0 {
			result["topic_planner"] = "" + strings.Join(topics, ", ")
		}
	} else {
		result["topic_planner"] = "No topic planner"
	}

	if lifeSkills != nil && len(lifeSkills.Options) > 0 {
		var skills []string
		for _, opt := range lifeSkills.Options {
			if opt.Name != "" {
				skills = append(skills, opt.Name)
			}
		}
		if len(skills) > 0 {
			result["life_skills"] = "" + strings.Join(skills, ", ")
		}
	} else {
		result["life_skills"] = "No life skills"
	}

	// Format Select Topic
	if selectTopic != nil && len(selectTopic.Options) > 0 {
		var topicNames []string
		for _, opt := range selectTopic.Options {
			if opt.TopicID != "" {
				topic, err := s.topicService.GetTopicInfor(ctx, opt.TopicID)
				if err != nil {
					return nil, fmt.Errorf("failed to get topic information: %w", err)
				}
				if topic != nil && topic.Name != "" {
					topicNames = append(topicNames, topic.Name)
				}
			}
		}
		if len(topicNames) > 0 {
			result["select_topic"] = "" + strings.Join(topicNames, ", ")
		}
	} else {
		result["select_topic"] = "No select topic"
	}

	return result, nil
}
