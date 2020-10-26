package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hatchery/components"
)

const (
	CollectionName = "numbers"
)

//Name Struct
type Name struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Used bool               `json:"used,omitempty" bson:"used,omitempty"`
}

func NewNameHatcheryMongo() *Name {
	return &Name{}
}

func (n *Name) Insert(db *components.StorageComponent) error {
	collection := db.GetCollection(CollectionName)

	id, err := collection.InsertOne(context.TODO(), n)
	if err != nil {
		return err
	}

	idValue, ok := id.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("Wrong id value")
	}

	n.ID = idValue

	return err
}

func (n *Name) Update(db *components.StorageComponent) error {
	collection := db.GetCollection(CollectionName)

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": n.ID}, n)
	if err != nil {
		return err
	}

	return err
}

func (n *Name) Delete(db *components.StorageComponent) error {
	collection := db.GetCollection(CollectionName)

	filter := bson.M{"name": n.Name}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return err
}

func NameByID(db *components.StorageComponent, id string) (*Name, error) {
	collection := db.GetCollection(CollectionName)

	filter := bson.M{"ID": id}
	var name *Name
	err := collection.FindOne(context.TODO(), filter).Decode(&name)
	if err != nil {
		return nil, err
	}

	return name, err
}

func GetNames(db *components.StorageComponent) ([]*Name, error) {

	collection := db.GetCollection(CollectionName)

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	// we created Book array
	var names []*Name
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var name *Name
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

func SafeGetName(db *components.StorageComponent) (*Name, error) {

	collection := db.GetCollection(CollectionName)

	var name *Name
	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"used": bson.M{"$exists": false}}

	err := collection.FindOne(context.TODO(), filter).Decode(&name)
	if err != nil {
		return nil, err
	}

	return name, err
}
