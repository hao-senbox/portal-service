package models

type RequestStudentActivity struct {
	StudentID    string      `json:"student_id" validate:"required"`
	Date         string      `json:"date" validate:"required"`
	TypeActivity string      `json:"type_activity" validate:"required"`
	Data         map[string]interface{} `json:"data" validate:"required"`
	AssignedBy   string      `json:"assigned_by" validate:"required"`
	SubmittedAt  string      `json:"submitted_at" validate:"required"`
}
