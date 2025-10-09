package ieb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IEB struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Owner       *Owner             `json:"owner" bson:"owner"`
	TermID      string             `json:"term_id" bson:"term_id"`
	LanguageID  int                `json:"language_id" bson:"language_id"`
	Information []Information      `json:"information" bson:"information"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Owner struct {
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	OwnerRole string `json:"owner_role" bson:"owner_role"`
}

type Information struct {
	Tilte    string    `json:"title" bson:"title"`
	Contents []Content `json:"contents" bson:"contents"`
}

type Content struct {
	Label     string    `json:"label" bson:"label"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
