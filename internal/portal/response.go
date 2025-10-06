package portal

import "time"

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
	SessionID    string                `json:"session_id,omitempty"`
	TypeActivity string                `json:"type_activity"`
	Data         []StudentActivityData `json:"data"`
	SubmittedAt  time.Time             `json:"submitted_at"`
	AssignedBy   string                `json:"assigned_by"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

