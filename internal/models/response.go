package models

import "time"

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
}

type StudentDailyActivities struct {
	StudentID  string            `bson:"student_id, omitempty" json:"student_id"`
	Date       string            `bson:"date" json:"date"`
	Activities []ActivitySummary `bson:"activities" json:"activities"`
}

type ActivitySummary struct {
	TypeActivity string              `bson:"type_activity" json:"type_activity"`
	Summary      ActivitySummaryData `json:"summary"`
	Details      []ActivityDetail    `json:"details"`
}

type ActivitySummaryData struct {
	TotalSessions int                    `json:"total_sessions"`
	Statistics    map[string]interface{} `json:"statistics"`
}

type ActivityDetail struct {
	SessionID   string                `json:"session_id,omitempty"`
	Data        []StudentActivityData `json:"data"`
	SubmittedAt time.Time             `json:"submitted_at"`
	AssignedBy  string                `json:"assigned_by"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

const (
	ErrInvalidOperation = "ERR_INVALID_OPERATION"
	ErrInvalidRequest   = "ERR_INVALID_REQUEST"
)
