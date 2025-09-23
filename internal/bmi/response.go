package bmi

import (
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BMIStudentResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Student   *user.UserInfor    `json:"student" bson:"student"`
	Date      string             `json:"date" bson:"date"`
	Height    float64            `json:"height" bson:"height"`
	Weight    float64            `json:"weight" bson:"weight"`
	Teacher   *user.UserInfor    `json:"teacher" bson:"teacher"`
	BMI       float64            `json:"bmi" bson:"bmi"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
