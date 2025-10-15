package teacherassign

type CreateTeacherAssignmentRequest struct {
	ParentID      string `json:"parent_id" bson:"parent_id"`
	StudentID     string `json:"student_id" bson:"student_id"`
	Month         int     `json:"month" bson:"month"`
	Year          int     `json:"year" bson:"year"`
	Language      *Data   `json:"language" bson:"language"`
	Qualification *Data   `bson:"qualification" json:"qualification"`
	ExperienceExt *Data   `bson:"experience_external" json:"experience_external"`
	ExperienceInt *Data   `bson:"experience_internal" json:"experience_internal"`
	PDLevel       *Data   `bson:"pd_level" json:"pd_level"`
	SkillSet      *Data   `bson:"skill_set" json:"skill_set"`
	AgeRange      *Data   `bson:"age_range" json:"age_range"`
}

type UpdateTeacherAssignmentRequest struct {
	Language      *Data   `json:"language" bson:"language"`
	Qualification *Data   `bson:"qualification" json:"qualification"`
	ExperienceExt *Data   `bson:"experience_external" json:"experience_external"`
	ExperienceInt *Data   `bson:"experience_internal" json:"experience_internal"`
	PDLevel       *Data   `bson:"pd_level" json:"pd_level"`
	SkillSet      *Data   `bson:"skill_set" json:"skill_set"`
	AgeRange      *Data   `bson:"age_range" json:"age_range"`
}
