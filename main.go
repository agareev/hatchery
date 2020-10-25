package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"hatchery/components"
	"hatchery/handlers"
	"net/http"
	"strconv"
)

func main() {
	cfg := components.NewConfiguration()

	//if err := cfg.ParseEnvManual(); err != nil {
	//	logrus.Fatal("Error upload conf from env: ", err)
	//}

	if err := cfg.ParseEnvPkg(); err != nil {
		logrus.Fatal("Error upload conf from env: ", err)
	}

	dbComponent := components.NewStorageComponent(cfg)
	err := dbComponent.ConnectDB()
	if err != nil {
		logrus.Fatal("Error connect to DB", err)
	}

	h := handlers.NewNumbersHandler(dbComponent)

	router := mux.NewRouter()
	// get all names
	router.HandleFunc("/", h.GetNames).Methods("GET")
	router.HandleFunc("/name", h.GetNames).Methods("GET")
	router.HandleFunc("/name/all", h.GetNames).Methods("GET")
	//
	router.HandleFunc("/name", h.CreateName).Methods("POST")
	router.HandleFunc("/name/{name}", h.DeleteName).Methods("DELETE") // unused
	// check next free name
	router.HandleFunc("/name/x", h.SafeGetName).Methods("GET")
	// use via zergrush
	router.HandleFunc("/name/x", h.XSetName).Methods("POST")
	router.HandleFunc("/name/loose", h.MakeFree).Methods("POST")

	logrus.Infof("Server is listening... on:%d", cfg.Rest.Port)
	logrus.Fatal(http.ListenAndServe(prepareAddr(cfg), router))
}

func prepareAddr(cfg *components.Configuration) string {
	return fmt.Sprintf(cfg.Rest.Host + ":" + strconv.Itoa(cfg.Rest.Port))
}
