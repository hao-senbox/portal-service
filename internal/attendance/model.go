package attendance

type AttendanceUserInfo struct {
	AttendanceID   string  `json:"id"`
	StudentID      string  `json:"student_id"`
	Date           string  `json:"date"`
	Temperature    float64 `json:"temperature"`
}
