package studypreference

type CreateStudyPreferenceRequest struct {
	OrganizationID   string `json:"organization_id" bson:"organization_id" binding:"required"`
	TermID           string `json:"term_id" bson:"term_id" binding:"required"`
	StudentID        string `json:"student_id" bson:"student_id" binding:"required"`
	ParentSelection  []Data `json:"parent_selection" bson:"parent_selection"`
	TeacherSelection []Data `json:"teacher_selection" bson:"teacher_selection" binding:"required"`
}

// UpdateStudyPreferenceRequest - Parent update their selection based on teacher's suggestion
type UpdateStudyPreferenceRequest struct {
	ParentSelections []Data `json:"parent_selections" bson:"parent_selections" binding:"required"`
}
