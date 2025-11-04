package selectoptions

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectOptions struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	TermID         string             `json:"term_id" bson:"term_id"`
	StudentID      string             `json:"student_id,omitempty" bson:"student_id,omitempty"`
	Type           string             `json:"type" bson:"type"` // "iep_priority", "topic_planner", "select_topic", "life_skills"
	Options        []Options          `json:"options" bson:"options"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	UpdatedBy      string             `json:"updated_by" bson:"updated_by"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	IsDeleted      bool               `json:"is_deleted" bson:"is_deleted"`
}

type Options struct {
	Name      string  `json:"name,omitempty" bson:"name,omitempty"`
	Order     int     `json:"order" bson:"order"`
	Icon      *string `json:"icon,omitempty" bson:"icon,omitempty"`         // Required for iep_priority, select_topic
	Status    *string `json:"status,omitempty" bson:"status,omitempty"`     // Required for topic_planner
	Optionals *string `json:"optional,omitempty" bson:"optional,omitempty"` // Required for topic_planner
	TopicID   string `json:"topic_id,omitempty" bson:"topic_id,omitempty"` // Required for topic_planner
}
