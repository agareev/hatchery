package components

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageComponent struct {
	DbName string
	Client *mongo.Client
}

func NewStorageComponent(cfg *Configuration) (*StorageComponent, error) {
	ctx := context.TODO()

	connectionStr := fmt.Sprintf("mongodb://%s:%s/", cfg.Mongo.Uri, cfg.Mongo.Port)
	clientOptions := options.Client().ApplyURI(connectionStr)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.WithError(err).Error("Error - mongo client connection")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("Error - mongo client ping")
		return nil, err
	}

	logrus.Info("Connected to MongoDB!")

	return &StorageComponent{
		Client: client,
		DbName: cfg.Mongo.Name,
	}, nil
}

func (s *StorageComponent) GetCollection(collectionName string) *mongo.Collection {
	return s.Client.Database(s.DbName).Collection(collectionName)
}
