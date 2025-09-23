package timer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Timer struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	StudentID       string             `json:"student_id" bson:"student_id"`
	StartColor      string             `json:"start_color" bson:"start_color"`
	EndColor        string             `json:"end_color" bson:"end_color"`
	Duration        int64              `json:"duration" bson:"duration"`
	CenterLine      int64              `json:"center_line" bson:"center_line"`
	NumberOfSound   int                `json:"number_of_sound" bson:"number_of_sound"`
	Image           string             `json:"image" bson:"image"`
	OpacityDuration float64            `json:"opacity_duration" bson:"opacity_duration"`
	TypePlay        string             `json:"type_play" bson:"type_play"`
	CreatedBy       string             `json:"created_by" bson:"created_by"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
