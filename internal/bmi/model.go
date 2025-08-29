package bmi

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BMI struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID string             `json:"student_id" bson:"student_id"`
	Date      time.Time             `json:"date" bson:"date"`
	Height    float64            `json:"height" bson:"height"`
	Weight    float64            `json:"weight" bson:"weight"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	BMI       float64            `json:"bmi" bson:"bmi"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
