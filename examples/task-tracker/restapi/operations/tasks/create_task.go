package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// CreateTaskHandlerFunc turns a function with the right signature into a create task handler
type CreateTaskHandlerFunc func(CreateTaskParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateTaskHandlerFunc) Handle(params CreateTaskParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CreateTaskHandler interface for that can handle valid create task params
type CreateTaskHandler interface {
	Handle(CreateTaskParams, interface{}) middleware.Responder
}

// NewCreateTask creates a new http.Handler for the create task operation
func NewCreateTask(ctx *middleware.ApiContext, handler CreateTaskHandler) *CreateTask {
	return &CreateTask{Context: ctx, Handler: handler}
}

/*CreateTask swagger:route POST /tasks tasks createTask

Creates a 'Task' object.

Allows for creating a task.
This operation requires authentication so that we know which user
created the task.


*/
type CreateTask struct {
	Context *middleware.ApiContext
	Params  CreateTaskParams
	Handler CreateTaskHandler
}

func (o *CreateTask) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewCreateTaskParams()

	uprinc, err := o.Context.Authorize(ctx, r, route)
	if err != nil {
		o.Context.Respond(ctx, rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(ctx, rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params, principal) // actually handle the request

	o.Context.Respond(ctx, rw, r, route.Produces, route, res)

}
