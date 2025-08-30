package timer

type Timer struct {
	ID            string `json:"id" bson:"_id"`
	StudentID     string `json:"student_id" bson:"student_id"`
	StartColor    string `json:"start_color" bson:"start_color"`
	EndColor      string `json:"end_color" bson:"end_color"`
	Duration      int64  `json:"duration" bson:"duration"`
	NumberOfSound int    `json:"number_of_sound" bson:"number_of_sound"`
	Image         string `json:"image" bson:"image"`
	TypePlay      string `json:"type_play" bson:"type_play"`
	CreatedBy     string `json:"created_by" bson:"created_by"`
	CreatedAt     string `json:"created_at" bson:"created_at"`
	UpdatedAt     string `json:"updated_at" bson:"updated_at"`
}
