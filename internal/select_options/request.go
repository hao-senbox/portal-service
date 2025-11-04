package selectoptions

type CreateSelectOptionRequest struct {
	OrganizationID string    `json:"organization_id" bson:"organization_id" binding:"required"`
	StudentID      string    `json:"student_id,omitempty" bson:"student_id,omitempty"`
	Type           string    `json:"type" bson:"type" binding:"required"`
	Options        []Options `json:"options" bson:"options" binding:"required"`
}

type UpdateSelectOptionRequest struct {
	Options *[]Options `json:"options,omitempty" bson:"options"`
}

