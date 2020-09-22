package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Name Struct
type Name struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Used bool               `json:"used,omitempty" bson:"used,omitempty"`
}
