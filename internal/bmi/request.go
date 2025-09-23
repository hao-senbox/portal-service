package bmi

type CreateBMIStudentRequest struct {
	StudentID string  `json:"student_id" bson:"student_id"`
	Date      string  `json:"date" bson:"date"`
	Height    float64 `json:"height" bson:"height"`
	Weight    float64 `json:"weight" bson:"weight"`
}
