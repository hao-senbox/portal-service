package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentActivity struct {
	ID           primitive.ObjectID    `bson:"_id, omitempty" json:"id"`
	StudentID    string                `bson:"student_id, omitempty" json:"student_id"`
	TypeActivity string                `bson:"type_activity" json:"type_activity"`
	Date         time.Time             `bson:"date" json:"date"`
	Data         []StudentActivityData `bson:"data" json:"data"`
	SubmittedAt  time.Time             `bson:"submitted_at" json:"submitted_at"`
	AssignedBy   string                `bson:"assigned_by" json:"assigned_by"`
	CreatedAt    time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at" json:"updated_at"`
}

type StudentActivityData struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}
