package drink

import (
	"portal/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrinkResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Student   *user.UserInfor    `json:"student" bson:"student"`
	Date      string             `json:"date" bson:"date"`
	Liquids   []Liquid           `json:"liquids" bson:"liquids"`
	Teacher   *user.UserInfor    `json:"teacher" bson:"teacher"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type DrinkDailyTotals struct {
	Date       string          `json:"date" bson:"date"`
	Student    *user.UserInfor `json:"student" bson:"student"`
	Statistics []Satistic      `json:"statistics" bson:"statistics"`
	Teacher    *user.UserInfor `json:"teacher" bson:"teacher"`
	CreatedAt  time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" bson:"updated_at"`
}

type Satistic struct {
	Type  string  `json:"type" bson:"type"`
	Total float64 `json:"total" bson:"total"`
}
