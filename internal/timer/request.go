package timer

type CreateTimerRequest struct {
	StudentID       string  `json:"student_id" bson:"student_id"`
	StartColor      string  `json:"start_color" bson:"start_color"`
	EndColor        string  `json:"end_color" bson:"end_color"`
	Duration        int64   `json:"duration" bson:"duration"`
	LineCenter      int64   `json:"line_center" bson:"line_center"`
	OpacityDuration float64 `json:"opacity_duration" bson:"opacity_duration"`
	NumberOfSound   int     `json:"number_of_sound" bson:"number_of_sound"`
	Image           string  `json:"image" bson:"image"`
	TypePlay        string  `json:"type_play" bson:"type_play"`
}
