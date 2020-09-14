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

func names(w http.ResponseWriter, r *http.Request) {
	log.Println("names")
	fmt.Fprint(w, "names")
}

func getAllNames(w http.ResponseWriter, r *http.Request) {
	log.Println("start requests")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println(err, "err with find")
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err, "err with cur")
		}
		log.Println("ssss", result)
	}
	if err := cur.Err(); err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, "names")
}

func init() {
	dburl = getEnv("DBURL", "localhost")
	dbport = getEnv("DBPORT", "27017")
	dbname = getEnv("DBNAME", "testing")
	httpport = getEnv("LISTEN_PORT", "8182")
	collectionName = getEnv("COLLECTIONNAME", "numbers")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dburl+":"+dbport))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to mongodb")

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", names)
	router.HandleFunc("/name", names)
	router.HandleFunc("/name/all", getAllNames)
	router.HandleFunc("/name/x", names)
	router.HandleFunc("/name/busy", names)
	router.HandleFunc("/name/loose", names)
	http.Handle("/", router)

	collection = client.Database(dbname).Collection(collectionName)

	log.Println("Server is listening... on :" + httpport)
	err := http.ListenAndServe(":"+httpport, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
