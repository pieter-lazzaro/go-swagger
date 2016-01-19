package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// FindHandlerFunc turns a function with the right signature into a find handler
type FindHandlerFunc func(FindParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn FindHandlerFunc) Handle(params FindParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// FindHandler interface for that can handle valid find params
type FindHandler interface {
	Handle(FindParams, interface{}) middleware.Responder
}

// NewFind creates a new http.Handler for the find operation
func NewFind(ctx *middleware.ApiContext, handler FindHandler) *Find {
	return &Find{Context: ctx, Handler: handler}
}

/*Find swagger:route GET / todos find

Find find API

*/
type Find struct {
	Context *middleware.ApiContext
	Params  FindParams
	Handler FindHandler
}

func (o *Find) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	o.Params = NewFindParams()

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
