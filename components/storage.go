package components

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hatchery/models"
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

	s.Client = client
	//collection := client.Database("runners_info").Collection("shared_name")
	logrus.Info("Connected to MongoDB!")
	return nil
}

func (s *StorageComponent) GetCollection() *mongo.Collection {
	return s.Client.Database(s.Config.Mongo.Name).Collection(s.Config.Mongo.CollectionName)
}

func (s *StorageComponent) GetNames() ([]*models.Name, error) {
	//Connection mongoDB with helper class
	collection := s.GetCollection()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	// we created Book array
	var names []*models.Name
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var name *models.Name
		// & character returns the memory address of the following variable.
		err := cur.Decode(&name) // decode similar to deserialize process.
		if err != nil {
			return nil, err
		}

		// add item our array
		names = append(names, name)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return names, err
}

func (s *StorageComponent) CreateName(model *models.Name) (*mongo.InsertOneResult, error) {
	collection := s.GetCollection()

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *StorageComponent) SafeGetName() (*models.Name, error) {

	//Connection mongoDB with helper class
	collection := s.GetCollection()

	var name *models.Name
	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"used": bson.M{"$exists": false}}

	err := collection.FindOne(context.TODO(), filter).Decode(&name)
	if err != nil {
		return nil, err
	}

	return name, err
}

func (s *StorageComponent) XSetName() (*models.Name, error) {
	//Connection mongoDB with helper class
	collection := s.GetCollection()

	// Create filter
	filter := bson.M{"used": bson.M{"$exists": false}}

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"used", true},
		}},
	}

	var name *models.Name
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&name)

	if err != nil {
		return nil, err
	}

	return name, err
}

func (s *StorageComponent) DeleteName(params map[string]string) (*mongo.DeleteResult, error) {
	// string to primitve.ObjectID
	// name, err := primitive.ObjectIDFromHex(params["name"])

	collection := s.GetCollection()

	// prepare filter.
	filter := bson.M{"name": params["name"]}

	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return deleteResult, err
}

func (s *StorageComponent) UpdateName(params map[string]string, name *models.Name) (*models.Name, error) {
	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := s.GetCollection()

	// Create filter
	filter := bson.M{"_id": id}

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", name.Name},
			{"user", name.Used},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&name)
	if err != nil {
		return nil, err
	}

	name.ID = id
	return name, err
}

func (s *StorageComponent) MakeFee(name *models.Name) error {
	//Connection mongoDB with helper class
	collection := s.GetCollection()
	// bson.M{}, we passed empty filter. So we want to get all data.
	filter := bson.M{"name": name.Name}
	update := bson.D{
		{"$unset", bson.D{
			{"used", ""},
		}},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
