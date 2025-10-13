package timer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Timer struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	StudentID         string             `json:"student_id" bson:"student_id"`
	StartColor        string             `json:"start_color" bson:"start_color"`
	EndColor          string             `json:"end_color" bson:"end_color"`
	Duration          int64              `json:"duration" bson:"duration"`
	CenterLine        int64              `json:"center_line" bson:"center_line"`
	NumberOfSound     int                `json:"number_of_sound" bson:"number_of_sound"`
	ImageStartKey     string             `json:"image_start_key" bson:"image_start_key"`
	ShowImageStart    bool               `json:"show_image_start" bson:"show_image_start"`
	CaptionImageStart string             `json:"caption_image_start" bson:"caption_image_start"`
	ImageEndKey       string             `json:"image_end_key" bson:"image_end_key"`
	ShowImageEnd      bool               `json:"show_image_end" bson:"show_image_end"`
	CaptionImageEnd   string             `json:"caption_image_end" bson:"caption_image_end"`
	OpacityDuration   float64            `json:"opacity_duration" bson:"opacity_duration"`
	TypePlay          string             `json:"type_play" bson:"type_play"`
	CreatedBy         string             `json:"created_by" bson:"created_by"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

type IsTime struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	StudentID         string             `json:"student_id" bson:"student_id"`
	IndexImage        int                `json:"index_image" bson:"index_image"`
	Sentence          string             `json:"sentence" bson:"sentence"`
	Mode              string             `json:"mode" bson:"mode"`
	ImageKey          string             `json:"image_key" bson:"image_key"`
	RewardImageKey    string             `json:"reward_image" bson:"reward_image"`
	BehaviourImageKey string             `json:"behaviour_image" bson:"behaviour_image"`
	CaptionImage      string             `json:"caption_image" bson:"caption_image"`
	ImageSize         string             `json:"image_size" bson:"image_size"`
	CreatedBy         string             `json:"created_by" bson:"created_by"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}
