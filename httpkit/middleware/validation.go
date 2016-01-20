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
	"log"
	"mime"
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/swag"

	"golang.org/x/net/context"
)

// NewValidation starts a new validation middleware
func newValidation(ctx *ApiContext, next Handler) Handler {

	return HandlerFunc(func(rCtx context.Context, rw http.ResponseWriter, r *http.Request) {
		matched := MatchedRouteFromContext(rCtx)

		_, result := ctx.BindAndValidate(rCtx, r, matched)

		if result != nil {
			ctx.Respond(rw, r, matched.Produces, matched, result)
			return
		}

		next.ServeHTTP(rCtx, rw, r)
	})
}

type validation struct {
	context *ApiContext
	result  []error
	request *http.Request
	route   *MatchedRoute
	bound   map[string]interface{}
}

type untypedBinder map[string]interface{}

func (ub untypedBinder) BindRequest(r *http.Request, route *MatchedRoute, consumer httpkit.Consumer) error {
	if err := route.Binder.Bind(r, route.Params, consumer, ub); err != nil {
		return err
	}
	return nil
}

// ContentType validates the content type of a request
func validateContentType(allowed []string, actual string) *errors.Validation {
	mt, _, err := mime.ParseMediaType(actual)
	if err != nil {
		return errors.InvalidContentType(actual, allowed)
	}
	if swag.ContainsStringsCI(allowed, mt) {
		return nil
	}
	return errors.InvalidContentType(actual, allowed)
}

func validateRequestContentType(ctx context.Context, route *MatchedRoute, r *http.Request) error {
	if !httpkit.CanHaveBody(r.Method) {
		return nil
	}

	ct := ContentTypeFromContext(ctx)

	if ct == nil {
		return errors.New(http.StatusBadRequest, "Could not read content type.")
	}

	if ct.Err != nil {
		log.Println("content type", ct)
		return ct.Err
	}

	if err := validateContentType(route.Consumes, ct.MediaType); err != nil {
		log.Println("content type", ct, err)
		return err
	}

	route.Consumer = route.Consumers[ct.MediaType]

	return nil
}

func validateRequestParameters(ctx context.Context, route *MatchedRoute, request *http.Request) (map[string]interface{}, []error) {
	var errs []error
	bound := make(map[string]interface{})

	if result := route.Binder.Bind(request, route.Params, route.Consumer, bound); result != nil {
		if result.Error() == "validation failure list" {
			for _, e := range result.(*errors.Validation).Value.([]interface{}) {
				errs = append(errs, e.(error))
			}
			return nil, errs
		}
	}

	return bound, errs
}

type boundParams struct {
	params map[string]interface{}
	errs   []error
}

func validateRequestContext(ctx context.Context, r *http.Request, route *MatchedRoute) boundParams {
	result := boundParams{
		params: make(map[string]interface{}),
	}

	if err := validateRequestContentType(ctx, route, r); err != nil {
		result.errs = append(result.errs, err)
	}

	if format := ResponseFormatFromContext(ctx); format == "" {
		result.errs = append(result.errs, errors.InvalidResponseFormat(r.Header.Get(httpkit.HeaderAccept), route.Produces))
	}

	if result.errs != nil {

		log.Printf("Bound Params w/ Errors: %#v\n", result)
		return result
	}

	bound, err := validateRequestParameters(ctx, route, r)

	if err != nil {
		result.errs = append(result.errs, err...)
	}

	result.params = bound

	log.Printf("Bound Params: %#v\n", result)
	return result
}

func validateRequest(ctx *ApiContext, request *http.Request, route *MatchedRoute) *validation {
	validate := &validation{
		context: ctx,
		request: request,
		route:   route,
		bound:   make(map[string]interface{}),
	}

	validate.contentType()
	validate.responseFormat()
	if len(validate.result) == 0 {
		validate.parameters()
	}

	return validate
}

func (v *validation) parameters() {
	if result := v.route.Binder.Bind(v.request, v.route.Params, v.route.Consumer, v.bound); result != nil {
		if result.Error() == "validation failure list" {
			for _, e := range result.(*errors.Validation).Value.([]interface{}) {
				v.result = append(v.result, e.(error))
			}
			return
		}
		v.result = append(v.result, result)
	}
}

func (v *validation) contentType() {
	if httpkit.CanHaveBody(v.request.Method) {
		ct, _, err := v.context.ContentType(v.request)
		if err != nil {
			v.result = append(v.result, err)
		} else {
			if err := validateContentType(v.route.Consumes, ct); err != nil {
				v.result = append(v.result, err)
			}
			v.route.Consumer = v.route.Consumers[ct]
		}
	}
}

func (v *validation) responseFormat() {
	if str := v.context.ResponseFormat(v.request, v.route.Produces); str == "" {
		v.result = append(v.result, errors.InvalidResponseFormat(v.request.Header.Get(httpkit.HeaderAccept), v.route.Produces))
	}
}
