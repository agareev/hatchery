// Code generated by go-swagger; DO NOT EDIT.

package name

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"hatchery/server/models"
)

// DeleteNameNameIDOKCode is the HTTP code returned for type DeleteNameNameIDOK
const DeleteNameNameIDOKCode int = 200

/*DeleteNameNameIDOK OK

swagger:response deleteNameNameIdOK
*/
type DeleteNameNameIDOK struct {
}

// NewDeleteNameNameIDOK creates DeleteNameNameIDOK with default headers values
func NewDeleteNameNameIDOK() *DeleteNameNameIDOK {

	return &DeleteNameNameIDOK{}
}

// WriteResponse to the client
func (o *DeleteNameNameIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteNameNameIDInternalServerErrorCode is the HTTP code returned for type DeleteNameNameIDInternalServerError
const DeleteNameNameIDInternalServerErrorCode int = 500

/*DeleteNameNameIDInternalServerError Not Ok

swagger:response deleteNameNameIdInternalServerError
*/
type DeleteNameNameIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.HatcheryError `json:"body,omitempty"`
}

// NewDeleteNameNameIDInternalServerError creates DeleteNameNameIDInternalServerError with default headers values
func NewDeleteNameNameIDInternalServerError() *DeleteNameNameIDInternalServerError {

	return &DeleteNameNameIDInternalServerError{}
}

// WithPayload adds the payload to the delete name name Id internal server error response
func (o *DeleteNameNameIDInternalServerError) WithPayload(payload *models.HatcheryError) *DeleteNameNameIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete name name Id internal server error response
func (o *DeleteNameNameIDInternalServerError) SetPayload(payload *models.HatcheryError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteNameNameIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}