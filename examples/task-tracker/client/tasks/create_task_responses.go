package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/strfmt"
)

// CreateTaskReader is a Reader for the CreateTask structure.
type CreateTaskReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the recieved o.
func (o *CreateTaskReader) ReadResponse(response client.Response, consumer httpkit.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewCreateTaskCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewCreateTaskDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	}
}

// NewCreateTaskCreated creates a CreateTaskCreated with default headers values
func NewCreateTaskCreated() *CreateTaskCreated {
	return &CreateTaskCreated{}
}

/*CreateTaskCreated handles this case with default header values.

Task created
*/
type CreateTaskCreated struct {
}

func (o *CreateTaskCreated) Error() string {
	return fmt.Sprintf("[POST /tasks][%d] createTaskCreated ", 201)
}

func (o *CreateTaskCreated) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateTaskDefault creates a CreateTaskDefault with default headers values
func NewCreateTaskDefault(code int) *CreateTaskDefault {
	return &CreateTaskDefault{
		_statusCode: code,
	}
}

/*CreateTaskDefault handles this case with default header values.

CreateTaskDefault create task default
*/
type CreateTaskDefault struct {
	_statusCode int
}

// Code gets the status code for the create task default response
func (o *CreateTaskDefault) Code() int {
	return o._statusCode
}

func (o *CreateTaskDefault) Error() string {
	return fmt.Sprintf("[POST /tasks][%d] createTask default ", o._statusCode)
}

func (o *CreateTaskDefault) readResponse(response client.Response, consumer httpkit.Consumer, formats strfmt.Registry) error {

	return nil
}
