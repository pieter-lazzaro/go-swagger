package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// UpdateTaskHandlerFunc turns a function with the right signature into a update task handler
type UpdateTaskHandlerFunc func(UpdateTaskParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateTaskHandlerFunc) Handle(params UpdateTaskParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// UpdateTaskHandler interface for that can handle valid update task params
type UpdateTaskHandler interface {
	Handle(UpdateTaskParams, interface{}) middleware.Responder
}

// NewUpdateTask creates a new http.Handler for the update task operation
func NewUpdateTask(ctx *middleware.ApiContext, handler UpdateTaskHandler) *UpdateTask {
	return &UpdateTask{Context: ctx, Handler: handler}
}

/*UpdateTask swagger:route PUT /tasks/{id} tasks updateTask

Updates the details for a task.

Allows for updating a task.
This operation requires authentication so that we know which user
last updated the task.


*/
type UpdateTask struct {
	Context *middleware.ApiContext
	Params  UpdateTaskParams
	Handler UpdateTaskHandler
}

func (o *UpdateTask) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewUpdateTaskParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
