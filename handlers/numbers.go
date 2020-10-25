package handlers

import (
	"encoding/json"
	"fmt"
	"hatchery/components"
	"hatchery/models"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	name, err := parseBody(r.Body)
	if err != nil {
		getError(err, w)
	}

	log.Println(name)

	result, err := n.DbComponent.CreateName(name)
	if err != nil {
		getError(err, w)
	}

	json.NewEncoder(w).Encode(result)
}

func (n *NumbersHandler) GetNames(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	names, err := n.DbComponent.GetNames()
	if err != nil {
		getError(err, w)
	}

	json.NewEncoder(w).Encode(names) // encode similar to serialize process.
}

func (n *NumbersHandler) SafeGetName(w http.ResponseWriter, r *http.Request) {
	name, err := n.DbComponent.SafeGetName()
	if err != nil {
		getError(err, w)
	}

	json.NewEncoder(w).Encode(name)
}

func (n *NumbersHandler) XSetName(w http.ResponseWriter, r *http.Request) {
	name, err := n.DbComponent.XSetName()
	if err != nil {
		getError(err, w)
	}

	fmt.Fprintf(w, name.Name)
}

func (n *NumbersHandler) DeleteName(w http.ResponseWriter, r *http.Request) {
	// Set header
	setupHeader(w)

	// get params
	var params = mux.Vars(r)

	deleteResult, err := n.DbComponent.DeleteName(params)
	if err != nil {
		getError(err, w)
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func (n *NumbersHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	var params = mux.Vars(r)

	name, err := parseBody(r.Body)
	if err != nil {
		getError(err, w)
	}

	newName, err := n.DbComponent.UpdateName(params, name)
	if err != nil {
		getError(err, w)
	}

	json.NewEncoder(w).Encode(newName)
}

func (n *NumbersHandler) MakeFree(w http.ResponseWriter, r *http.Request) {
	setupHeader(w)

	name, err := parseBody(r.Body)
	if err != nil {
		getError(err, w)
	}

	log.Println(name.Name)

	err = n.DbComponent.MakeFee(name)
	if err != nil {
		getError(err, w)
	}

	fmt.Fprintf(w, name.Name)
}

func parseBody(body io.ReadCloser) (*models.Name, error) {
	var name *models.Name
	// we decode our body request params
	err := json.NewDecoder(body).Decode(&name)
	if err != nil {
		return nil, err
	}

	return name, err
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
