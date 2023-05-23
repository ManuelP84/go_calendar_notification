package task

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoSettings "github.com/ManuelP84/calendar_notification/infra/mongo"
)

const (
	dbName         = "tasks"
	collectionName = "taskEvents"
)

type MongoRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoRepository(settings *mongoSettings.MongoDbSettings) *MongoRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	addr := fmt.Sprintf("mongodb://%s:%s", settings.Host, settings.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))

	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)

	col := db.Collection(collectionName)

	return &MongoRepository{db, col}
}

func (mongo *MongoRepository) InsertEvent(ctx context.Context, eventType string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := mongo.collection.InsertOne(ctxTimeout, bson.D{{Key: "event", Value: eventType}})

	if err != nil {
		return err
	}

	//id := res.InsertedID.(primitive.ObjectID).String()

	return nil
}
