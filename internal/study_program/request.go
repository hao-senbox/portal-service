package studyprogram

type CreateStudyProgramRequest struct {
	ParentID      string  `json:"parent_id" bson:"parent_id"`
	StudentID     string  `json:"student_id" bson:"student_id"`
	Month         int     `json:"month" bson:"month"`
	Year          int     `json:"year" bson:"year"`
	TimeSlot      *Data   `json:"time_slot" bson:"time_slot"`
	ServiceRatio  *Data   `json:"service_ratio" bson:"service_ratio"`
	SkillPercent  *Data   `json:"skill_percent" bson:"skill_percent"`
	TeacherWeight *Data   `json:"teacher_weight" bson:"teacher_weight"`
	Extras        []*Data `json:"extras" bson:"extras"`
	OtherFees     []*Data `json:"other_fees" bson:"other_fees"`
}
type UpdateStudyProgramRequest struct {
	Month         int     `json:"month" bson:"month"`
	Year          int     `json:"year" bson:"year"`
	TimeSlot      *Data   `json:"time_slot" bson:"time_slot"`
	ServiceRatio  *Data   `json:"service_ratio" bson:"service_ratio"`
	SkillPercent  *Data   `json:"skill_percent" bson:"skill_percent"`
	TeacherWeight *Data   `json:"teacher_weight" bson:"teacher_weight"`
	Extras        []*Data `json:"extras" bson:"extras"`
	OtherFees     []*Data `json:"other_fees" bson:"other_fees"`
}
