package models

import (
	"time"
)

type StudentActivity struct {
	StudentID    string                 `bson:"student_id, omitempty" json:"student_id"`
	TypeActivity string                 `bson:"type_activity" json:"type_activity"`
	Date         time.Time              `bson:"date" json:"date"`
	Data         map[string]interface{} `bson:"data" json:"data"`
	SubittedAt   time.Time              `bson:"submitted_at" json:"submitted_at"`
	CreatdBy     string                 `bson:"created_by" json:"created_by"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}
