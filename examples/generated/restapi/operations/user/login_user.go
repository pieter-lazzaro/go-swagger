package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"golang.org/x/net/context"
)

// LoginUserHandlerFunc turns a function with the right signature into a login user handler
type LoginUserHandlerFunc func(LoginUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn LoginUserHandlerFunc) Handle(params LoginUserParams) middleware.Responder {
	return fn(params)
}

// LoginUserHandler interface for that can handle valid login user params
type LoginUserHandler interface {
	Handle(LoginUserParams) middleware.Responder
}

// NewLoginUser creates a new http.Handler for the login user operation
func NewLoginUser(ctx *middleware.ApiContext, handler LoginUserHandler) *LoginUser {
	return &LoginUser{Context: ctx, Handler: handler}
}

/*LoginUser swagger:route GET /user/login user loginUser

Logs user into the system

*/
type LoginUser struct {
	Context *middleware.ApiContext
	Params  LoginUserParams
	Handler LoginUserHandler
}

func (o *LoginUser) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewLoginUserParams()

	if err := o.Context.BindValidRequest(r, route, &o.Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(o.Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
