package task

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
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

func (mongo *MongoRepository) InsertEvent(ctx context.Context, event events.TaskEvent) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := mongo.collection.InsertOne(ctxTimeout, bson.D{
		{Key: "type", Value: event.EventType},
		{Key: "taskID", Value: event.Task.Id},
		{Key: "title", Value: event.Task.Title},
		{Key: "description", Value: event.Task.Description},
	})

	if err != nil {
		return err
	}

	//id := res.InsertedID.(primitive.ObjectID).String()

	return nil
}
