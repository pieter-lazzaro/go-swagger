package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// GetTaskCommentsHandlerFunc turns a function with the right signature into a get task comments handler
type GetTaskCommentsHandlerFunc func(GetTaskCommentsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTaskCommentsHandlerFunc) Handle(params GetTaskCommentsParams) middleware.Responder {
	return fn(params)
}

// GetTaskCommentsHandler interface for that can handle valid get task comments params
type GetTaskCommentsHandler interface {
	Handle(GetTaskCommentsParams) middleware.Responder
}

// NewGetTaskComments creates a new http.Handler for the get task comments operation
func NewGetTaskComments(ctx *middleware.ApiContext, handler GetTaskCommentsHandler) *GetTaskComments {
	return &GetTaskComments{Context: ctx, Handler: handler}
}

/*GetTaskComments swagger:route GET /tasks/{id}/comments tasks getTaskComments

Gets the comments for a task

The comments require a size parameter.


*/
type GetTaskComments struct {
	Context *middleware.ApiContext
	Params  GetTaskCommentsParams
	Handler GetTaskCommentsHandler
}

func (o *GetTaskComments) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	o.Params = NewGetTaskCommentsParams()

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
