// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"net/http"
	"strings"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware/untyped"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/gorilla/context"
	netContext "golang.org/x/net/context"
)

// A Builder can create middlewares
type Builder func(Handler) Handler


type HandlerFunc func(netContext.Context, http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(ctx netContext.Context, w http.ResponseWriter, r *http.Request) {
    h(ctx, w, r)
}

type Handler interface {
    ServeHTTP(netContext.Context, http.ResponseWriter, *http.Request)
}

// PassthroughBuilder returns the handler, aka the builder identity function
func PassthroughBuilder(handler Handler) Handler { return handler }

// RequestBinder is an interface for types to implement
// when they want to be able to bind from a request
type RequestBinder interface {
	BindRequest(*http.Request, *MatchedRoute) error
}

// Responder is an interface for types to implement
// when they want to be considered for writing HTTP responses
type Responder interface {
	WriteResponse(http.ResponseWriter, httpkit.Producer)
}

// Context is a type safe wrapper around an untyped request context
// used throughout to store request context with the gorilla context module
type ApiContext struct {
	spec    *spec.Document
	api     RoutableAPI
	router  Router
	formats strfmt.Registry
}

type routableUntypedAPI struct {
	api             *untyped.API
	handlers        map[string]map[string]Handler
	defaultConsumes string
	defaultProduces string
}

func newRoutableUntypedAPI(spec *spec.Document, api *untyped.API, context *ApiContext) *routableUntypedAPI {
	var handlers map[string]map[string]Handler
	if spec == nil || api == nil {
		return nil
	}
	for method, hls := range spec.Operations() {
		um := strings.ToUpper(method)
		for path, op := range hls {
			schemes := spec.SecurityDefinitionsFor(op)

			if oh, ok := api.OperationHandlerFor(method, path); ok {
				if handlers == nil {
					handlers = make(map[string]map[string]Handler)
				}
				if b, ok := handlers[um]; !ok || b == nil {
					handlers[um] = make(map[string]Handler)
				}

				handlers[um][path] = HandlerFunc(func(rCtx netContext.Context, w http.ResponseWriter, r *http.Request) {
					// lookup route info in the context
					route := MatchedRouteFromContext(rCtx)

					// bind and validate the request using reflection
					bound, validation := context.BindAndValidate(rCtx, r, route)
					if validation != nil {
						context.Respond(w, r, route.Produces, route, validation)
						return
					}

					// actually handle the request
					result, err := oh.Handle(bound)
					if err != nil {
						// respond with failure
						context.Respond(w, r, route.Produces, route, err)
						return
					}

					// respond with success
					context.Respond(w, r, route.Produces, route, result)
				})

				if len(schemes) > 0 {
					handlers[um][path] = newSecureAPI(context, handlers[um][path])
				}
			}
		}
	}

	return &routableUntypedAPI{
		api:             api,
		handlers:        handlers,
		defaultProduces: api.DefaultProduces,
		defaultConsumes: api.DefaultConsumes,
	}
}

func (r *routableUntypedAPI) HandlerFor(method, path string) (Handler, bool) {
	paths, ok := r.handlers[strings.ToUpper(method)]
	if !ok {
		return nil, false
	}
	handler, ok := paths[path]
	return handler, ok
}
func (r *routableUntypedAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return r.api.ServeError
}
func (r *routableUntypedAPI) ConsumersFor(mediaTypes []string) map[string]httpkit.Consumer {
	return r.api.ConsumersFor(mediaTypes)
}
func (r *routableUntypedAPI) ProducersFor(mediaTypes []string) map[string]httpkit.Producer {
	return r.api.ProducersFor(mediaTypes)
}
func (r *routableUntypedAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]httpkit.Authenticator {
	return r.api.AuthenticatorsFor(schemes)
}
func (r *routableUntypedAPI) Formats() strfmt.Registry {
	return r.api.Formats()
}

func (r *routableUntypedAPI) DefaultProduces() string {
	return r.defaultProduces
}

func (r *routableUntypedAPI) DefaultConsumes() string {
	return r.defaultConsumes
}

// NewRoutableContext creates a new context for a routable API
func NewRoutableContext(spec *spec.Document, routableAPI RoutableAPI, routes Router) *ApiContext {
	ctx := &ApiContext{spec: spec, api: routableAPI}
	return ctx
}

// NewContext creates a new context wrapper
func NewContext(spec *spec.Document, api *untyped.API, routes Router) *ApiContext {
	ctx := &ApiContext{spec: spec}
	ctx.api = newRoutableUntypedAPI(spec, api, ctx)
	return ctx
}

// Serve serves the specified spec with the specified api registrations as a http.Handler
func Serve(spec *spec.Document, api *untyped.API) http.Handler {
	return ServeWithBuilder(spec, api, PassthroughBuilder)
}

// ServeWithBuilder serves the specified spec with the specified api registrations as a http.Handler that is decorated
// by the Builder
func ServeWithBuilder(spec *spec.Document, api *untyped.API, builder Builder) http.Handler {
	context := NewContext(spec, api, nil)
	return context.APIHandler(builder)
}

type contextKey int8

const (
	_ contextKey = iota
	ctxContentType
	ctxResponseFormat
	ctxMatchedRoute
	ctxAllowedMethods
	ctxBoundParams
	ctxSecurityPrincipal

	ctxConsumer
)

type contentTypeValue struct {
	MediaType string
	Charset   string
	Err       error
}

// BasePath returns the base path for this API
func (c *ApiContext) BasePath() string {
	return c.spec.BasePath()
}

// RequiredProduces returns the accepted content types for responses
func (c *ApiContext) RequiredProduces() []string {
	return c.spec.RequiredProduces()
}

// BindValidRequest binds a params object to a request but only when the request is valid
// if the request is not valid an error will be returned
func (c *ApiContext) BindValidRequest(request *http.Request, route *MatchedRoute, binder RequestBinder) error {
	var res []error

	// check and validate content type, select consumer
	if httpkit.CanHaveBody(request.Method) {
		ct, _, err := httpkit.ContentType(request.Header)
		if err != nil {
			res = append(res, err)
		} else {
			if err := validateContentType(route.Consumes, ct); err != nil {
				res = append(res, err)
			}
			route.Consumer = route.Consumers[ct]
		}
	}

	// check and validate the response format
	if len(res) == 0 {
		if str := NegotiateContentType(request, route.Produces, ""); str == "" {
			res = append(res, errors.InvalidResponseFormat(request.Header.Get(httpkit.HeaderAccept), route.Produces))
		}
	}

	// now bind the request with the provided binder
	// it's assumed the binder will also validate the request and return an error if the
	// request is invalid
	if binder != nil && len(res) == 0 {
		if err := binder.BindRequest(request, route); err != nil {
			res = append(res, err)
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func NewContextWithContentType(ctx netContext.Context, request *http.Request) netContext.Context {

	mt, cs, err := httpkit.ContentType(request.Header)
	if err != nil {
		return netContext.WithValue(ctx, ctxContentType, &contentTypeValue{"", "", err})
	}

	return netContext.WithValue(ctx, ctxContentType, &contentTypeValue{mt, cs, nil})
}

func ContentTypeFromContext(ctx netContext.Context) *contentTypeValue {
	if v, ok := ctx.Value(ctxContentType).(*contentTypeValue); ok {
		return v
	}

	return nil
}

// ContentType gets the parsed value of a content type
func (c *ApiContext) ContentType(request *http.Request) (string, string, *errors.ParseError) {
	if v, ok := context.GetOk(request, ctxContentType); ok {
		if val, ok := v.(*contentTypeValue); ok {
			return val.MediaType, val.Charset, nil
		}
	}

	mt, cs, err := httpkit.ContentType(request.Header)
	if err != nil {
		return "", "", err
	}
	context.Set(request, ctxContentType, &contentTypeValue{mt, cs, nil})
	return mt, cs, nil
}

// LookupRoute looks a route up and returns true when it is found
func (c *ApiContext) LookupRoute(request *http.Request) (*MatchedRoute, bool) {
	if route, ok := c.router.Lookup(request.Method, request.URL.Path); ok {
		return route, ok
	}
	return nil, false
}

func NewContextWithMatchedRoute(ctx netContext.Context, route *MatchedRoute) netContext.Context {
	return netContext.WithValue(ctx, ctxMatchedRoute, route)
}

func MatchedRouteFromContext(ctx netContext.Context) *MatchedRoute {
	if v, ok := ctx.Value(ctxMatchedRoute).(*MatchedRoute); ok {
		return v
	}

	return nil
}

func (c *ApiContext) NewRequestContext(r *http.Request) netContext.Context {
	ctx := netContext.TODO()

	ctx = NewContextWithContentType(ctx, r)

	if route, ok := c.LookupRoute(r); ok {

		ctx = NewContextWithMatchedRoute(ctx, route)

		ctx = NewContextWithSecurityPrincipal(ctx, r, route)

		ctx = NewContextWithResponseFormat(ctx, route, r)

		ctx = NewContextWithBoundParams(ctx, r, route)
	}

	return ctx
}

// RouteInfo tries to match a route for this request
func (c *ApiContext) RouteInfo(request *http.Request) (*MatchedRoute, bool) {
	if v, ok := context.GetOk(request, ctxMatchedRoute); ok {
		if val, ok := v.(*MatchedRoute); ok {
			return val, ok
		}
	}

	if route, ok := c.LookupRoute(request); ok {
		context.Set(request, ctxMatchedRoute, route)
		return route, ok
	}

	return nil, false
}

func NewContextWithResponseFormat(ctx netContext.Context, route *MatchedRoute, request *http.Request) netContext.Context {
	if format := NegotiateContentType(request, route.Produces, ""); format != "" {
		return netContext.WithValue(ctx, ctxResponseFormat, format)
	}
	return ctx
}

func ResponseFormatFromContext(ctx netContext.Context) string {
	if v, ok := ctx.Value(ctxResponseFormat).(string); ok {
		return v
	}

	return ""
}

// ResponseFormat negotiates the response content type
func (c *ApiContext) ResponseFormat(r *http.Request, offers []string) string {
	if v, ok := context.GetOk(r, ctxResponseFormat); ok {
		if val, ok := v.(string); ok {
			return val
		}
	}

	format := NegotiateContentType(r, offers, "")
	if format != "" {
		context.Set(r, ctxResponseFormat, format)
	}
	return format
}

// AllowedMethods gets the allowed methods for the path of this request
func (c *ApiContext) AllowedMethods(request *http.Request) []string {
	return c.router.OtherMethods(request.Method, request.URL.Path)
}

func NewContextWithSecurityPrincipal(ctx netContext.Context, request *http.Request, route *MatchedRoute) netContext.Context {

	if len(route.Authenticators) == 0 {
		return ctx
	}

	for _, authenticator := range route.Authenticators {
		applies, usr, err := authenticator.Authenticate(request)

		if !applies || err != nil || usr == nil {
			continue
		}
		return netContext.WithValue(ctx, ctxSecurityPrincipal, usr)
	}

	return ctx
}

func SecurityPrincipalFromContext(ctx netContext.Context) interface{} {
	return ctx.Value(ctxSecurityPrincipal)
}

// Authorize authorizes the request byt checking if a security principal is needed and has been set
func (c *ApiContext) Authorize(ctx netContext.Context, request *http.Request, route *MatchedRoute) (interface{}, error) {

	// No auth needed
	if len(route.Authenticators) == 0 {
		return nil, nil
	}

	principal := SecurityPrincipalFromContext(ctx)

	if principal == nil {
		return nil, errors.Unauthenticated("invalid credentials")
	}

	return principal, nil
}

func NewContextWithBoundParams(ctx netContext.Context, request *http.Request, route *MatchedRoute) netContext.Context {
	params := validateRequestContext(ctx, request, route)
	return netContext.WithValue(ctx, ctxBoundParams, params)
}

func BoundParamsFromContext(ctx netContext.Context) boundParams {
	if v, ok := ctx.Value(ctxBoundParams).(boundParams); ok {
		return v
	}

	return boundParams{errs: []error{errors.New(http.StatusBadRequest, "Could not read parameters.")}}
}

// BindAndValidate binds and validates the request
func (c *ApiContext) BindAndValidate(ctx netContext.Context, request *http.Request, matched *MatchedRoute) (interface{}, error) {
	params := BoundParamsFromContext(ctx)

	if len(params.errs) > 0 {
		return params.params, errors.CompositeValidationError(params.errs...)
	}

	return params.params, nil
}

// NotFound the default not found responder for when no route has been matched yet
func (c *ApiContext) NotFound(rw http.ResponseWriter, r *http.Request) {
	c.Respond(rw, r, []string{c.api.DefaultProduces()}, nil, errors.NotFound("not found"))
}

// Respond renders the response after doing some content negotiation
func (c *ApiContext) Respond(rw http.ResponseWriter, r *http.Request, produces []string, route *MatchedRoute, data interface{}) {
	offers := []string{c.api.DefaultProduces()}
	for _, mt := range produces {
		if mt != c.api.DefaultProduces() {
			offers = append(offers, mt)
		}
	}

	format := c.ResponseFormat(r, offers)
	rw.Header().Set(httpkit.HeaderContentType, format)

	if resp, ok := data.(Responder); ok {
		producers := route.Producers
		prod, ok := producers[format]
		if !ok {
			panic(errors.New(http.StatusInternalServerError, "can't find a producer for "+format))
		}
		resp.WriteResponse(rw, prod)
		return
	}

	if err, ok := data.(error); ok {
		if format == "" {
			rw.Header().Set(httpkit.HeaderContentType, httpkit.JSONMime)
		}
		if route == nil || route.Operation == nil {
			c.api.ServeErrorFor("")(rw, r, err)
			return
		}
		c.api.ServeErrorFor(route.Operation.ID)(rw, r, err)
		return
	}

	if route == nil || route.Operation == nil {
		rw.WriteHeader(200)
		if r.Method == "HEAD" {
			return
		}
		producers := c.api.ProducersFor(offers)
		prod, ok := producers[format]
		if !ok {
			panic(errors.New(http.StatusInternalServerError, "can't find a producer for "+format))
		}
		if err := prod.Produce(rw, data); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
		return
	}

	if _, code, ok := route.Operation.SuccessResponse(); ok {
		rw.WriteHeader(code)
		if code == 204 || r.Method == "HEAD" {
			return
		}

		producers := route.Producers
		prod, ok := producers[format]
		if !ok {
			panic(errors.New(http.StatusInternalServerError, "can't find a producer for "+format))
		}
		if err := prod.Produce(rw, data); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
		return
	}

	c.api.ServeErrorFor(route.Operation.ID)(rw, r, errors.New(http.StatusInternalServerError, "can't produce response"))
}

// APIHandler returns a handler to serve
func (c *ApiContext) APIHandler(builder Builder) http.Handler {
	b := builder
	if b == nil {
		b = PassthroughBuilder
	}
	return specMiddleware(c, newRouter(c, b(newOperationExecutor(c))))
}
