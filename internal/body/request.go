package body

type CreateCheckInRequest struct {
	StudentID string `json:"student_id" bson:"student_id"`
	Date      string `json:"date" bson:"date"`
	Context   string `json:"context" bson:"context"`
	Gender    *string `json:"gender" bson:"gender"`
	Type      string `json:"type" bson:"type"`
	Marks     []Mark `json:"marks" bson:"marks"`
	CreatedBy string `json:"created_by" bson:"created_by"`
}