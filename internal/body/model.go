package body

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckIn struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID string             `json:"student_id" bson:"student_id"`
	Date      time.Time          `json:"date" bson:"date"`
	Context   string             `json:"context" bson:"context"` // home, school
	Gender    *string            `json:"gender" bson:"gender"`   // male, female
	Type      string             `json:"type" bson:"type"`       // dressed_front, dressed_back || body_front, body_back || face, feeling
	Marks     []Mark             `json:"marks" bson:"marks"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Mark struct {
	Name        string    `json:"name" bson:"name"`
	Note        *string   `json:"note" bson:"note"`
	Severity    int       `json:"severity" bson:"severity"`
	SubmittedAt time.Time `json:"submitted_at" bson:"submitted_at"`
}
