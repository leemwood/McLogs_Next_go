package storage

import (
	"context"
	"mclogs-go/internal/config"
	"mclogs-go/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client     *mongo.Client
	collection *mongo.Collection
	cfg        *config.StorageConfig
}

func NewMongoStorage(cfg *config.Config) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoDB.URL))
	if err != nil {
		return nil, err
	}

	coll := client.Database(cfg.Database.MongoDB.DB).Collection(cfg.Database.MongoDB.Collection)

	return &MongoStorage{
		client:     client,
		collection: coll,
		cfg:        &cfg.Storage,
	}, nil
}

func (s *MongoStorage) Put(ctx context.Context, content string) (string, error) {
	rawID := GenerateRawID()
	id := GetFullID(s.cfg.CurrentID, rawID)
	log := models.Log{
		ID:        id,
		Content:   content,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.TTL) * time.Second),
	}

	_, err := s.collection.InsertOne(ctx, log)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *MongoStorage) Get(ctx context.Context, id string) (*models.Log, error) {
	var log models.Log
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (s *MongoStorage) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *MongoStorage) Renew(ctx context.Context, id string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"expires_at": time.Now().Add(time.Duration(s.cfg.TTL) * time.Second)}},
	)
	return err
}
