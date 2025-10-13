package program_planner

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgramPlaner struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	StudentID      string             `json:"student_id" bson:"student_id"`
	Month          int                `json:"month" bson:"month"`
	Year           int                `json:"year" bson:"year"`
	TotalFee       float64            `json:"total_fee" bson:"total_fee"`
	SelectedSlots  []SelectedSlot     `json:"selected_slots" bson:"selected_slots"`
	Weeks          []WeekPlan         `json:"weeks" bson:"weeks"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	IsDeleted      bool               `json:"is_deleted" bson:"is_deleted"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type SelectedSlot struct {
	TimeRange string   `json:"time_range" bson:"time_range"`
	Days      []string `json:"days" bson:"days"`
	Selected  bool     `json:"selected" bson:"selected"`
	Fee       float64  `json:"fee" bson:"fee"`
}

type WeekPlan struct {
	WeekNumber int         `json:"week_number" bson:"week_number"`
	WeekStart  time.Time   `json:"week_start" bson:"week_start"`
	WeekEnd    time.Time   `json:"week_end" bson:"week_end"`
	WeekFee    float64     `json:"week_fee" bson:"week_fee"`
	Slots      []DailySlot `json:"slots" bson:"slots"`
}

type DailySlot struct {
	DayOfWeek string  `json:"day_of_week" bson:"day_of_week"`
	Time      string  `json:"time" bson:"time"`
	Selected  bool    `json:"selected" bson:"selected"`
	Fee       float64 `json:"fee" bson:"fee"`
}
