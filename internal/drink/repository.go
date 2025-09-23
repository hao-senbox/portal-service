package drink

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DrinkRepository interface {
	CreateDrink(ctx context.Context, drink *Drink) (string, error)
	GetDrinks(ctx context.Context, studentID string, date *time.Time) ([]*Drink, error)
	GetDrink(ctx context.Context, id primitive.ObjectID) (*Drink, error)
}

type drinkRepository struct {
	collection *mongo.Collection
}

func NewDrinkRepository(collection *mongo.Collection) DrinkRepository {
	return &drinkRepository{
		collection: collection,
	}
}

func (d *drinkRepository) CreateDrink(ctx context.Context, drink *Drink) (string, error) {

	result, err := d.collection.InsertOne(ctx, drink)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (d *drinkRepository) GetDrinks(ctx context.Context, studentID string, date *time.Time) ([]*Drink, error) {

	filter := bson.M{}

	if studentID != "" {
		filter["student_id"] = studentID
	}

	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour)

		filter["date"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}
	findOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var drinks []*Drink

	cursor, err := d.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &drinks)
	if err != nil {
		return nil, err
	}

	return drinks, nil

}

func (d *drinkRepository) GetDrink(ctx context.Context, id primitive.ObjectID) (*Drink, error) {

	var drink Drink

	err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&drink)
	if err != nil {
		return nil, err
	}

	return &drink, nil

}
