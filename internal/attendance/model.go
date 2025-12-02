package attendance

type AttendanceUserInfo struct {
	AttendanceID string  `json:"id"`
	StudentID    string  `json:"student_id"`
	CheckInTime  string  `json:"check_in_time"`
	CheckOutTime string  `json:"check_out_time"`
	Date         string  `json:"date"`
	Temperature  float64 `json:"temperature"`
}
