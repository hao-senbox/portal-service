package teacherassign

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherAssignment struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ParentID      string             `json:"parent_id" bson:"parent_id"`
	StudentID     string             `json:"student_id" bson:"student_id"`
	Month         int                `json:"month" bson:"month"`
	Year          int                `json:"year" bson:"year"`
	Language      *Data              `json:"language" bson:"language"`
	Qualification *Data              `bson:"qualification" json:"qualification"`
	ExperienceExt *Data              `bson:"experience_external" json:"experience_external"`
	ExperienceInt *Data              `bson:"experience_internal" json:"experience_internal"`
	PDLevel       *Data              `bson:"pd_level" json:"pd_level"`
	SkillSet      *Data              `bson:"skill_set" json:"skill_set"`
	AgeRange      *Data              `bson:"age_range" json:"age_range"`
	MonthlyFee    float64            `bson:"monthly_fee" json:"monthly_fee"`
	CreatedBy     string             `bson:"created_by" json:"created_by"`
	IsDeleted     bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Data struct {
	Label string  `json:"label"`
	Price float64 `json:"price"`
}
