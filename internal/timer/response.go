package timer

import (
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimerResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Student         *user.UserInfor    `json:"student" bson:"student"`
	StartColor      string             `json:"start_color" bson:"start_color"`
	EndColor        string             `json:"end_color" bson:"end_color"`
	Duration        int64              `json:"duration" bson:"duration"`
	CenterLine      int64              `json:"center_line" bson:"center_line"`
	OpacityDuration float64            `json:"opacity_duration" bson:"opacity_duration"`
	NumberOfSound   int                `json:"number_of_sound" bson:"number_of_sound"`
	Image           string             `json:"image" bson:"image"`
	TypePlay        string             `json:"type_play" bson:"type_play"`
	Teacher         *user.UserInfor    `json:"teacher" bson:"teacher"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}
