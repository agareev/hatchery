package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"hatchery/helper"
	"hatchery/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection                                      *mongo.Collection
	client                                          *mongo.Client
	ctx                                             = context.TODO()
	dburl, dbport, dbname, collectionName, httpport string
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var names []models.Name

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
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

func safeGetName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var name models.Name

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"used": bson.M{"$exists": false}}
	err := collection.FindOne(context.TODO(), filter).Decode(&name)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(name)
}

func xsetName(w http.ResponseWriter, r *http.Request) {
	var name models.Name

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()

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
		helper.GetError(err, w)
		return
	}

	fmt.Fprintf(w, name.Name)

}
func createName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var name models.Name

	// we decode our body request params
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		log.Println(err)
	}

	// connect db
	collection := helper.ConnectDB()

	log.Println(name)
	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), name)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func deleteName(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	// name, err := primitive.ObjectIDFromHex(params["name"])

	collection := helper.ConnectDB()

	// prepare filter.
	filter := bson.M{"name": params["name"]}

	deleteResult, err := collection.DeleteMany(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func updateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var name models.Name

	collection := helper.ConnectDB()

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
		helper.GetError(err, w)
		return
	}

	name.ID = id

	json.NewEncoder(w).Encode(name)
}

func makeFree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var name models.Name

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()
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

func init() {
	dburl = getEnv("DBURL", "localhost")
	dbport = getEnv("DBPORT", "27017")
	dbname = getEnv("DBNAME", "testing")
	httpport = getEnv("LISTEN_PORT", "8182")
	collectionName = getEnv("COLLECTIONNAME", "numbers")
}

func main() {

	router := mux.NewRouter()
	// get all names
	router.HandleFunc("/", getNames).Methods("GET")
	router.HandleFunc("/name", getNames).Methods("GET")
	router.HandleFunc("/name/all", getNames).Methods("GET")
	//
	router.HandleFunc("/name", createName).Methods("POST")
	router.HandleFunc("/name/{name}", deleteName).Methods("DELETE") // unused
	// check next free name
	router.HandleFunc("/name/x", safeGetName).Methods("GET")
	// use via zergrush
	router.HandleFunc("/name/x", xsetName).Methods("POST")

	router.HandleFunc("/name/loose", makeFree).Methods("POST")

	log.Println("Server is listening... on :8182")
	log.Fatal(http.ListenAndServe(":8182", router))
}
