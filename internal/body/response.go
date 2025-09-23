package body

import (
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckInReponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Student   *user.UserInfor    `json:"student" bson:"student"`
	Date      string          `json:"date" bson:"date"`
	Context   string             `json:"context" bson:"context"`
	Gender    *string            `json:"gender" bson:"gender"`
	Type      string             `json:"type" bson:"type"`
	Marks     []Mark             `json:"marks" bson:"marks"`
	Teacher   *user.UserInfor    `json:"teacher" bson:"teacher"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}