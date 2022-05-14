package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repo[T any] interface {
	Disconnect() error
	ReadObject(id string) (*T, error)
	SaveObject(object *T) (string, error)
}

type MongoRepo[T any] struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func Connect[T any](uri, db, collectionName string) (Repo[T], error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	database := client.Database(db)
	collection := database.Collection(collectionName)
	return &MongoRepo[T]{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (m *MongoRepo[T]) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (m *MongoRepo[T]) ReadObject(id string) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := m.collection.Find(ctx, primitive.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	var t T
	err = cursor.Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m *MongoRepo[T]) SaveObject(object *T) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := m.collection.InsertOne(ctx, object)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
