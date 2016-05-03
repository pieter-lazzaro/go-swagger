package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/go-swagger/go-swagger/examples/task-tracker/models"
)

/*UploadTaskFileCreated File added

swagger:response uploadTaskFileCreated
*/
type UploadTaskFileCreated struct {
}

// NewUploadTaskFileCreated creates UploadTaskFileCreated with default headers values
func NewUploadTaskFileCreated() *UploadTaskFileCreated {
	return &UploadTaskFileCreated{}
}

// WriteResponse to the client
func (o *UploadTaskFileCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
}

/*UploadTaskFileDefault Error response

swagger:response uploadTaskFileDefault
*/
type UploadTaskFileDefault struct {
	_statusCode int
	/*
	  Required: true
	*/
	XErrorCode string `json:"X-Error-Code"`

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewUploadTaskFileDefault creates UploadTaskFileDefault with default headers values
func NewUploadTaskFileDefault(code int) *UploadTaskFileDefault {
	if code <= 0 {
		code = 500
	}

	return &UploadTaskFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the upload task file default response
func (o *UploadTaskFileDefault) WithStatusCode(code int) *UploadTaskFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the upload task file default response
func (o *UploadTaskFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithXErrorCode adds the xErrorCode to the upload task file default response
func (o *UploadTaskFileDefault) WithXErrorCode(xErrorCode string) *UploadTaskFileDefault {
	o.XErrorCode = xErrorCode
	return o
}

// SetXErrorCode sets the xErrorCode to the upload task file default response
func (o *UploadTaskFileDefault) SetXErrorCode(xErrorCode string) {
	o.XErrorCode = xErrorCode
}

// WithPayload adds the payload to the upload task file default response
func (o *UploadTaskFileDefault) WithPayload(payload *models.Error) *UploadTaskFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload task file default response
func (o *UploadTaskFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadTaskFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header X-Error-Code
	rw.Header().Add("X-Error-Code", fmt.Sprintf("%v", o.XErrorCode))

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
