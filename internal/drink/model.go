package drink

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Drink struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date      time.Time          `json:"date" bson:"date"`
	StudentID string             `json:"student_id" bson:"student_id"`
	Liquids   []Liquid           `json:"liquids" bson:"liquids"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Liquid struct {
	Type      string  `json:"type" bson:"type"`
	Amount    float64 `json:"amount" bson:"amount"`
	Capacity  float64 `json:"capacity" bson:"capacity"`
	Initial   float64 `json:"initial" bson:"initial"`
	Remaining float64 `json:"remaining" bson:"remaining"`
}
