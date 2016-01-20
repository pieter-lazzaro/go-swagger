package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// CreateUsersWithListInputHandlerFunc turns a function with the right signature into a create users with list input handler
type CreateUsersWithListInputHandlerFunc func(CreateUsersWithListInputParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateUsersWithListInputHandlerFunc) Handle(params CreateUsersWithListInputParams) middleware.Responder {
	return fn(params)
}

// CreateUsersWithListInputHandler interface for that can handle valid create users with list input params
type CreateUsersWithListInputHandler interface {
	Handle(CreateUsersWithListInputParams) middleware.Responder
}

// NewCreateUsersWithListInput creates a new http.Handler for the create users with list input operation
func NewCreateUsersWithListInput(ctx *middleware.ApiContext, handler CreateUsersWithListInputHandler) *CreateUsersWithListInput {
	return &CreateUsersWithListInput{Context: ctx, Handler: handler}
}

/*CreateUsersWithListInput swagger:route POST /user/createWithList user createUsersWithListInput

Creates list of users with given input array

*/
type CreateUsersWithListInput struct {
	Context *middleware.ApiContext
	Params  CreateUsersWithListInputParams
	Handler CreateUsersWithListInputHandler
}

func (o *CreateUsersWithListInput) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewCreateUsersWithListInputParams()

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(ctx, rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params) // actually handle the request

	o.Context.Respond(ctx, rw, r, route.Produces, route, res)

}
