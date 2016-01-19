package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/go-swagger/go-swagger/httpkit/validate"
	"github.com/go-swagger/go-swagger/strfmt"
	"golang.org/x/net/context"
)

// AddCommentToTaskHandlerFunc turns a function with the right signature into a add comment to task handler
type AddCommentToTaskHandlerFunc func(AddCommentToTaskParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn AddCommentToTaskHandlerFunc) Handle(params AddCommentToTaskParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// AddCommentToTaskHandler interface for that can handle valid add comment to task params
type AddCommentToTaskHandler interface {
	Handle(AddCommentToTaskParams, interface{}) middleware.Responder
}

// NewAddCommentToTask creates a new http.Handler for the add comment to task operation
func NewAddCommentToTask(ctx *middleware.ApiContext, handler AddCommentToTaskHandler) *AddCommentToTask {
	return &AddCommentToTask{Context: ctx, Handler: handler}
}

/*AddCommentToTask swagger:route POST /tasks/{id}/comments tasks addCommentToTask

Adds a comment to a task

The comment can contain ___github markdown___ syntax.
Fenced codeblocks etc are supported through pygments.


*/
type AddCommentToTask struct {
	Context *middleware.ApiContext
	Params  AddCommentToTaskParams
	Handler AddCommentToTaskHandler
}

func (o *AddCommentToTask) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	route := middleware.MatchedRouteFromContext(ctx)
	o.Params = NewAddCommentToTaskParams()

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

/*AddCommentToTaskBody A comment to create

These values can have github flavored markdown.


swagger:model AddCommentToTaskBody
*/
type AddCommentToTaskBody struct {

	/* Content content

	Required: true
	*/
	Content string `json:"content,omitempty"`

	/* UserID user id

	Required: true
	*/
	UserID int64 `json:"userId,omitempty"`
}

// Validate validates this add comment to task body
func (o *AddCommentToTaskBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateContent(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := o.validateUserID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AddCommentToTaskBody) validateContent(formats strfmt.Registry) error {

	if err := validate.RequiredString("body"+"."+"content", "body", string(o.Content)); err != nil {
		return err
	}

	return nil
}

func (o *AddCommentToTaskBody) validateUserID(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"userId", "body", int64(o.UserID)); err != nil {
		return err
	}

	return nil
}
