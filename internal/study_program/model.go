package studyprogram

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyProgram struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	ParentID      string             `json:"parent_id" bson:"parent_id"`
	StudentID     string             `json:"student_id" bson:"student_id"`
	Month         int                `json:"month" bson:"month"`
	Year          int                `json:"year" bson:"year"`
	TimeSlot      *Data              `bson:"time_slot" json:"time_slot"`         // M-F 8-12
	ServiceRatio  *Data              `bson:"service_ratio" json:"service_ratio"` // 1:1, 1:2
	SkillPercent  *Data              `bson:"skill_percent" json:"skill_percent"`
	TeacherWeight *Data              `bson:"teacher_weight" json:"teacher_weight"`
	Extras        []*Data            `bson:"extras" json:"extras"`
	OtherFees     []*Data            `bson:"other_fees" json:"other_fees"`
	MonthlyTotal  float64            `bson:"monthly_total" json:"monthly_total"`
	CreatedBy     string             `bson:"created_by" json:"created_by"`
	IsDeleted     bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

type Data struct {
	Label string  `json:"label"`
	Price float64 `json:"price"`
}
