package pet

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// FindPetsByStatusHandlerFunc turns a function with the right signature into a find pets by status handler
type FindPetsByStatusHandlerFunc func(FindPetsByStatusParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn FindPetsByStatusHandlerFunc) Handle(params FindPetsByStatusParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// FindPetsByStatusHandler interface for that can handle valid find pets by status params
type FindPetsByStatusHandler interface {
	Handle(FindPetsByStatusParams, interface{}) middleware.Responder
}

// NewFindPetsByStatus creates a new http.Handler for the find pets by status operation
func NewFindPetsByStatus(ctx *middleware.ApiContext, handler FindPetsByStatusHandler) *FindPetsByStatus {
	return &FindPetsByStatus{Context: ctx, Handler: handler}
}

/*FindPetsByStatus swagger:route GET /pet/findByStatus pet findPetsByStatus

Finds Pets by status

Multiple status values can be provided with comma seperated strings

*/
type FindPetsByStatus struct {
	Context *middleware.ApiContext
	Params  FindPetsByStatusParams
	Handler FindPetsByStatusHandler
}

func (o *FindPetsByStatus) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewFindPetsByStatusParams()

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
