package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hatchery/components"
	"hatchery/models"
	"log"
	"net/http"
)

type NumbersHandler struct {
	DbComponent *components.StorageComponent
}

func NewNumbersHandler(db *components.StorageComponent) *NumbersHandler {
	return &NumbersHandler{
		DbComponent: db,
	}
}

func (n *NumbersHandler) CreateName(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)
	var name models.Name

	// we decode our body request params
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		log.Println(err)
	}

	log.Println(name)
	collection := n.DbComponent.GetCollection()
	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), name)
	if err != nil {
		getError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (n *NumbersHandler) GetNames(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	// we created Book array
	var names []models.Name

	//Connection mongoDB with helper class
	collection := n.DbComponent.GetCollection()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		getError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var name models.Name
		// & character returns the memory address of the following variable.
		err := cur.Decode(&name) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		names = append(names, name)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(names) // encode similar to serialize process.
}

func (n *NumbersHandler) SafeGetName(w http.ResponseWriter, r *http.Request) {
	var name models.Name

	//Connection mongoDB with helper class
	collection := n.DbComponent.GetCollection()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"used": bson.M{"$exists": false}}
	err := collection.FindOne(context.TODO(), filter).Decode(&name)

	if err != nil {
		getError(err, w)
		return
	}

	json.NewEncoder(w).Encode(name)
}

func (n *NumbersHandler) XSetName(w http.ResponseWriter, r *http.Request) {
	var name models.Name

	//Connection mongoDB with helper class
	collection := n.DbComponent.GetCollection()

	// Create filter
	filter := bson.M{"used": bson.M{"$exists": false}}

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"used", true},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&name)

	if err != nil {
		getError(err, w)
		return
	}

	fmt.Fprintf(w, name.Name)
}

func (n *NumbersHandler) DeleteName(w http.ResponseWriter, r *http.Request) {
	// Set header
	setupHeader(w)

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	// name, err := primitive.ObjectIDFromHex(params["name"])

	collection := n.DbComponent.GetCollection()

	// prepare filter.
	filter := bson.M{"name": params["name"]}

	deleteResult, err := collection.DeleteMany(context.TODO(), filter)

	if err != nil {
		getError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func (n *NumbersHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var name models.Name

	collection := n.DbComponent.GetCollection()

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&name)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", name.Name},
			{"user", name.Used},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&name)

	if err != nil {
		getError(err, w)
		return
	}

	name.ID = id

	json.NewEncoder(w).Encode(name)
}

func (n *NumbersHandler) MakeFree(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	var name models.Name

	//Connection mongoDB with helper class
	collection := n.DbComponent.GetCollection()
	// we decode our body request params
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		log.Println(err)
	}

	log.Println(name.Name)
	// bson.M{}, we passed empty filter. So we want to get all data.
	filter := bson.M{"name": name.Name}
	update := bson.D{
		{"$unset", bson.D{
			{"used", ""},
		}},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, name.Name)
}

// getError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func getError(err error, w http.ResponseWriter) {

	log.Println(err.Error())

	response := models.ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
