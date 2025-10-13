package program_planner

type CreateProgramPlanerRequest struct {
	StudentID      string         `json:"student_id" bson:"student_id"`
	OrganizationID string         `json:"organization_id" bson:"organization_id"`
	Month          int            `json:"month" bson:"month"`
	Year           int            `json:"year" bson:"year"`
	TotalFee       float64        `json:"total_fee" bson:"total_fee"`
	SelectedSlots  []SelectedSlot `json:"selected_slots" bson:"selected_slots"`
}

type UpdateProgramPlanerRequest struct {
	Month          int            `json:"month" bson:"month"`
	Year           int            `json:"year" bson:"year"`
	SelectedSlots  []SelectedSlot `json:"selected_slots" bson:"selected_slots"`
}