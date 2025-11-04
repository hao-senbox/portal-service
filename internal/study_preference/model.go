package studypreference

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyPreference struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	OrganizationID   string             `json:"organization_id" bson:"organization_id"`
	TermID           string             `json:"term_id" bson:"term_id"`
	StudentID        string             `json:"student_id,omitempty" bson:"student_id,omitempty"`
	ParentSelection  []Data             `json:"parent_selections" bson:"parent_selections"`
	TeacherSelection []Data             `json:"teacher_selections" bson:"teacher_selections"`
	CreatedBy        string             `json:"created_by" bson:"created_by"`
	UpdatedBy        string             `json:"updated_by" bson:"updated_by"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	IsDeleted        bool               `json:"is_deleted" bson:"is_deleted"`
}

type Data struct {
	Pairs []Pair `json:"pairs" bson:"pairs"`
	Selected string	 `json:"selected" bson:"selected"`
}

type Pair struct {
	Category string `json:"category" bson:"category"`
	Value    int `json:"value" bson:"value"`
}
