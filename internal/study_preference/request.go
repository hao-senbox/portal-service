package studypreference

type CreateStudyPreferenceRequest struct {
	OrganizationID   string `json:"organization_id" bson:"organization_id" binding:"required"`
	StudentID        string `json:"student_id" bson:"student_id" binding:"required"`
	ParentSelection  []Data `json:"parent_selection" bson:"parent_selection"`
	TeacherSelection []Data `json:"teacher_selections" bson:"teacher_selections" binding:"required"`
}

// UpdateStudyPreferenceRequest - Parent update their selection based on teacher's suggestion
type UpdateStudyPreferenceRequest struct {
	ParentSelections []Data `json:"parent_selections" bson:"parent_selections" binding:"required"`
}
