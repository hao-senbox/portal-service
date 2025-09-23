package drink

type CreateDrinkRequest struct {
	StudentID string   `json:"student_id" bson:"student_id"`
	Date      string   `json:"date" bson:"date"`
	Liquids   []Liquid `json:"liquids" bson:"liquids"`
}
