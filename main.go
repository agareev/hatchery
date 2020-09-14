package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection                            *mongo.Collection
	dburl, dbport, dbname, collectionName string
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func names(w http.ResponseWriter, r *http.Request) {
	log.Println("names")
	fmt.Fprint(w, "names")
}

func init() {
	dburl = getEnv("DBURL", "localhost")
	dbport = getEnv("DBPORT", "27017")
	dbname = getEnv("DBNAME", "testing")
	collectionName = getEnv("COLLECTIONNAME", "numbers")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", names)
	router.HandleFunc("/name", names)
	router.HandleFunc("/name/all", names)
	router.HandleFunc("/name/x", names)
	router.HandleFunc("/name/busy", names)
	router.HandleFunc("/name/loose", names)
	http.Handle("/", router)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dburl+":"+dbport))

	collection = client.Database(dbname).Collection(collectionName)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	id := res.InsertedID
	log.Println(id)

	fmt.Println("Server is listening... on :8182")
	err = http.ListenAndServe(":8182", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
