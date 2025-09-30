package timer

type CreateTimerRequest struct {
	StudentID         string  `json:"student_id" bson:"student_id"`
	StartColor        string  `json:"start_color" bson:"start_color"`
	EndColor          string  `json:"end_color" bson:"end_color"`
	Duration          int64   `json:"duration" bson:"duration"`
	LineCenter        int64   `json:"line_center" bson:"line_center"`
	OpacityDuration   float64 `json:"opacity_duration" bson:"opacity_duration"`
	NumberOfSound     int     `json:"number_of_sound" bson:"number_of_sound"`
	ImageStartKey     string  `json:"image_start_key" bson:"image_start_key"`
	ShowImageStart    bool    `json:"show_image_start" bson:"show_image_start"`
	CaptionImageStart string  `json:"caption_image_start" bson:"caption_image_start"`
	ImageEndKey       string  `json:"image_end_key" bson:"image_end_key"`
	ShowImageEnd      bool    `json:"show_image_end" bson:"show_image_end"`
	CaptionImageEnd   string  `json:"caption_image_end" bson:"caption_image_end"`
	TypePlay          string  `json:"type_play" bson:"type_play"`
}

type CreateIsTimeRequest struct {
	StudentID    string `json:"student_id" bson:"student_id"`
	Sentence     string `json:"sentence" bson:"sentence"`
	Mode         string `json:"mode" bson:"mode"`
	ImageKey     string `json:"image_key" bson:"image_key"`
	CaptionImage string `json:"caption_image" bson:"caption_image"`
	ImageSize    string `json:"image_size" bson:"image_size"`
	CreatedBy    string `json:"created_by" bson:"created_by"`
}
