package ieb

type CreateIEBRequest struct {
	Owner       *Owner        `json:"owner" bson:"owner"`
	TermID      string        `json:"term_id" bson:"term_id"`
	LanguageID  string        `json:"language_id" bson:"language_id"`
	Information []Information `json:"information" bson:"information"`
}
