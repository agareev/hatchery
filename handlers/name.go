package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"hatchery/components"
	"hatchery/models"
	api "hatchery/server/models"
	"hatchery/server/restapi/operations/name"

	"net/http"
)

type NamesHandler struct {
	DbComponent *components.StorageComponent
}

func NewNamesHandler(db *components.StorageComponent) *NamesHandler {
	return &NamesHandler{
		DbComponent: db,
	}
}

func (h *NamesHandler) CreateName(params name.PostNameParams) middleware.Responder {
	if params.Payload.Name == nil {
		payload := prepareErrorResponse(http.StatusInternalServerError, fmt.Sprint("Error - incoming name is nil"))
		return name.NewPostNameInternalServerError().WithPayload(payload)
	}

	incoming := models.NewNameHatcheryMongo()
	incoming.Name = *params.Payload.Name

	err := incoming.Insert(h.DbComponent)
	if err != nil {
		payload := prepareErrorResponse(http.StatusInternalServerError, fmt.Sprint("Error - mongo creating new instance"))
		return name.NewPostNameInternalServerError().WithPayload(payload)
	}

	payload := convertDbModelToSwagger(incoming)
	return name.NewPostNameOK().WithPayload(payload)
}

func (h *NamesHandler) GetNames(_ name.GetParams) middleware.Responder {

	names, err := models.GetNames(h.DbComponent)
	if err != nil {
		payload := prepareErrorResponse(http.StatusInternalServerError, fmt.Sprint("Error - mongo creating new instance"))
		return name.NewPostNameInternalServerError().WithPayload(payload)
	}

	payload := api.NameCollectionResponse{}
	for _, nameFromDb := range names {
		payload = append(payload, convertDbModelToSwagger(nameFromDb))
	}

	return name.NewGetNameAllOK().WithPayload(payload)
}

func (h *NamesHandler) SafeGetName() {
	//TODO
}

func (h *NamesHandler) XSetName() {
	//TODO
}

func (h *NamesHandler) DeleteName(params name.DeleteNameNameIDParams) middleware.Responder {
	//TODO
	return nil
}

func (h *NamesHandler) UpdateName() {
	//TODO
}

func (h *NamesHandler) MakeFree() {
	//TODO
}

func prepareErrorResponse(code int64, message string) *api.HatcheryError {
	return &api.HatcheryError{
		Code:   code,
		Detail: message,
	}
}

func convertDbModelToSwagger(dbModelName *models.Name) *api.NameResponse {
	id := dbModelName.ID.String()
	return &api.NameResponse{
		ID:   &id,
		Name: &dbModelName.Name,
	}
}
