package ieb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEBRepository interface {
	CreateIEB(ctx context.Context, data *IEB) error
	GetIEB(ctx context.Context, userID string, termID string, languageID int64) (*IEB, error)
}

type iebRepository struct {
	IEBCollection *mongo.Collection
}

func NewIEBRepository(collection *mongo.Collection) IEBRepository {
	return &iebRepository{
		IEBCollection: collection,
	}
}

func (repository *iebRepository) CreateIEB(ctx context.Context, data *IEB) error {

	filter := bson.M{
		"owner.owner_id": data.Owner.OwnerID,
		"term_id":        data.TermID,
		"language_id":    data.LanguageID,
	}

	_, err := repository.IEBCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	_, err = repository.IEBCollection.InsertOne(ctx, data)
	return err

}

func (repository *iebRepository) GetIEB(ctx context.Context, userID string, termID string, languageID int64) (*IEB, error) {

	filter := bson.M{
		"owner.owner_id": userID,
		"term_id":        termID,
		"language_id":    languageID,
	}

	var ieb IEB

	err := repository.IEBCollection.FindOne(ctx, filter).Decode(&ieb)

	return &ieb, err
	
}
