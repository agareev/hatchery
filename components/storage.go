package components

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageComponent struct {
	Config *Configuration
	Client *mongo.Client
}

func NewStorageComponent(cfg *Configuration) *StorageComponent {
	return &StorageComponent{
		Config: cfg,
	}
}

func (s *StorageComponent) ConnectDB() error {
	ctx := context.TODO()

	connectionStr := fmt.Sprintf("mongodb://%s:%s/", s.Config.Mongo.Uri, s.Config.Mongo.Port)
	clientOptions := options.Client().ApplyURI(connectionStr)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.WithError(err).Error("Error - mongo client connection")
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("Error - mongo client ping")
		return err
	}

	//collection := client.Database("runners_info").Collection("shared_name")
	logrus.Info("Connected to MongoDB!")
	return nil
}

func (s *StorageComponent) GetCollection() *mongo.Collection {
	return s.Client.Database(s.Config.Mongo.Name).Collection(s.Config.Mongo.CollectionName)
}
