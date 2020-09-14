package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func names(w http.ResponseWriter, r *http.Request) {
	log.Println("names")
	fmt.Fprint(w, "names")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", names)
	http.Handle("/", router)

	fmt.Println("Server is listening... on :8182")
	err := http.ListenAndServe(":8182", nil)
	if err != nil {
		log.Fatalln(err)
	}

}
