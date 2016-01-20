package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// DestroyOneHandlerFunc turns a function with the right signature into a destroy one handler
type DestroyOneHandlerFunc func(DestroyOneParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DestroyOneHandlerFunc) Handle(params DestroyOneParams) middleware.Responder {
	return fn(params)
}

// DestroyOneHandler interface for that can handle valid destroy one params
type DestroyOneHandler interface {
	Handle(DestroyOneParams) middleware.Responder
}

// NewDestroyOne creates a new http.Handler for the destroy one operation
func NewDestroyOne(ctx *middleware.ApiContext, handler DestroyOneHandler) *DestroyOne {
	return &DestroyOne{Context: ctx, Handler: handler}
}

/*DestroyOne swagger:route DELETE /{id} todos destroyOne

DestroyOne destroy one API

*/
type DestroyOne struct {
	Context *middleware.ApiContext
	Params  DestroyOneParams
	Handler DestroyOneHandler
}

func (o *DestroyOne) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewDestroyOneParams()

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(ctx, rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params) // actually handle the request

	o.Context.Respond(ctx, rw, r, route.Produces, route, res)

}
