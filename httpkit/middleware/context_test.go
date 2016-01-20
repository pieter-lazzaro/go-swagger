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
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/internal/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	handler := Serve(spec, api)

	// serve spec document
	request, _ := http.NewRequest("GET", "http://localhost:8080/swagger.json", nil)
	request.Header.Add("Content-Type", httpkit.JSONMime)
	request.Header.Add("Accept", httpkit.JSONMime)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	request, _ = http.NewRequest("GET", "http://localhost:8080/swagger-ui", nil)
	recorder = httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
	assert.Equal(t, 404, recorder.Code)
}

func TestContextAuthorize(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	ctx := NewContext(spec, api, nil)
	ctx.router = DefaultRouter(spec, ctx.api)

	request, _ := httpkit.JSONRequest("GET", "/pets", nil)

	rctx := ctx.NewRequestContext(request)

	ri, ok := ctx.LookupRoute(request)
	assert.True(t, ok)

	p, err := ctx.Authorize(rctx, request, ri)
	assert.Error(t, err)
	assert.Nil(t, p)

	request.SetBasicAuth("wrong", "wrong")
	rctx = ctx.NewRequestContext(request)

	p, err = ctx.Authorize(rctx, request, ri)
	assert.Error(t, err)
	assert.Nil(t, p)

	request.SetBasicAuth("admin", "admin")
	rctx = ctx.NewRequestContext(request)
	p, err = ctx.Authorize(rctx, request, ri)

	assert.NoError(t, err)
	assert.Equal(t, "admin", p)
}

func TestContextBindAndValidate(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	ctx := NewContext(spec, api, nil)
	ctx.router = DefaultRouter(spec, ctx.api)

	request, _ := http.NewRequest("POST", "/pets", nil)
	request.Header.Add("Accept", "*/*")
	request.Header.Add("content-type", "text/html")

	ri, _ := ctx.LookupRoute(request)

	rctx := ctx.NewRequestContext(request)

	data, result := ctx.BindAndValidate(rctx, request, ri) // this requires a much more thorough test
	assert.NotNil(t, data)
	assert.NotNil(t, result)

	dd, rr := ctx.BindAndValidate(rctx, request, ri)
	assert.Equal(t, data, dd)
	assert.Equal(t, result, rr)
}

func TestContextRender(t *testing.T) {
	ct := httpkit.JSONMime
	spec, api := petstore.NewAPI(t)

	assert.NotNil(t, spec)
	assert.NotNil(t, api)
	ctx := NewContext(spec, api, nil)
	ctx.router = DefaultRouter(spec, ctx.api)

	request, _ := http.NewRequest("GET", "pets", nil)
	request.Header.Set(httpkit.HeaderAccept, ct)
	ri, _ := ctx.LookupRoute(request)
	log.Println(ri)

	rctx := ctx.NewRequestContext(request)
	recorder := httptest.NewRecorder()
	ctx.Respond(rctx, recorder, request, []string{ct}, ri, map[string]interface{}{"name": "hello"})
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "{\"name\":\"hello\"}\n", recorder.Body.String())

	recorder = httptest.NewRecorder()
	rctx = ctx.NewRequestContext(request)
	ctx.Respond(rctx, recorder, request, []string{ct}, ri, errors.New("this went wrong"))
	assert.Equal(t, 500, recorder.Code)

	recorder = httptest.NewRecorder()
	rctx = ctx.NewRequestContext(request)
	assert.Panics(t, func() { ctx.Respond(rctx, recorder, request, []string{ct}, ri, map[int]interface{}{1: "hello"}) })

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "pets", nil)
	rctx = ctx.NewRequestContext(request)
	assert.Panics(t, func() { ctx.Respond(rctx, recorder, request, []string{}, ri, map[string]interface{}{"name": "hello"}) })

	request, _ = http.NewRequest("GET", "/pets", nil)
	request.Header.Set(httpkit.HeaderAccept, ct)
	ri, _ = ctx.LookupRoute(request)

	recorder = httptest.NewRecorder()
	rctx = ctx.NewRequestContext(request)
	ctx.Respond(rctx, recorder, request, []string{ct}, ri, map[string]interface{}{"name": "hello"})
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "{\"name\":\"hello\"}\n", recorder.Body.String())

	recorder = httptest.NewRecorder()
	rctx = ctx.NewRequestContext(request)
	ctx.Respond(rctx, recorder, request, []string{ct}, ri, errors.New("this went wrong"))
	assert.Equal(t, 500, recorder.Code)

	recorder = httptest.NewRecorder()
	rctx = ctx.NewRequestContext(request)
	assert.Panics(t, func() { ctx.Respond(rctx, recorder, request, []string{ct}, ri, map[int]interface{}{1: "hello"}) })

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)
	rctx = ctx.NewRequestContext(request)
	assert.Panics(t, func() { ctx.Respond(rctx, recorder, request, []string{}, ri, map[string]interface{}{"name": "hello"}) })

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/pets/1", nil)
	ri, _ = ctx.LookupRoute(request)
	rctx = ctx.NewRequestContext(request)
	ctx.Respond(rctx, recorder, request, ri.Produces, ri, nil)
	assert.Equal(t, 204, recorder.Code)
}

// func TestContextValidResponseFormat(t *testing.T) {
// 	ct := "application/json"
// 	spec, api := petstore.NewAPI(t)
// 	ctx := NewContext(spec, api, nil)
// 	ctx.router = DefaultRouter(spec, ctx.api)

// 	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
// 	request.Header.Set(httpkit.HeaderAccept, ct)

// 	// check there's nothing there
// 	cached, ok := context.GetOk(request, ctxResponseFormat)
// 	assert.False(t, ok)
// 	assert.Empty(t, cached)

// 	// trigger the parse
// 	mt := ctx.ResponseFormat(request, []string{ct})
// 	assert.Equal(t, ct, mt)

// 	// check it was cached
// 	cached, ok = context.GetOk(request, ctxResponseFormat)
// 	assert.True(t, ok)
// 	assert.Equal(t, ct, cached)

// 	// check if the cast works and fetch from cache too
// 	mt = ctx.ResponseFormat(request, []string{ct})
// 	assert.Equal(t, ct, mt)
// }

// func TestContextInvalidResponseFormat(t *testing.T) {
// 	ct := "application/x-yaml"
// 	other := "application/sgml"
// 	spec, api := petstore.NewAPI(t)
// 	ctx := NewContext(spec, api, nil)
// 	ctx.router = DefaultRouter(spec, ctx.api)

// 	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
// 	request.Header.Set(httpkit.HeaderAccept, ct)

// 	// check there's nothing there
// 	cached, ok := context.GetOk(request, ctxResponseFormat)
// 	assert.False(t, ok)
// 	assert.Empty(t, cached)

// 	// trigger the parse
// 	mt := ctx.ResponseFormat(request, []string{other})
// 	assert.Empty(t, mt)

// 	// check it was cached
// 	cached, ok = context.GetOk(request, ctxResponseFormat)
// 	assert.False(t, ok)
// 	assert.Empty(t, cached)

// 	// check if the cast works and fetch from cache too
// 	mt = ctx.ResponseFormat(request, []string{other})
// 	assert.Empty(t, mt)
// }

// func TestContextValidRoute(t *testing.T) {
// 	spec, api := petstore.NewAPI(t)
// 	ctx := NewContext(spec, api, nil)
// 	ctx.router = DefaultRouter(spec, ctx.api)

// 	request, _ := http.NewRequest("GET", "/pets", nil)

// 	// check there's nothing there
// 	_, ok := context.GetOk(request, ctxMatchedRoute)
// 	assert.False(t, ok)

// 	matched, ok := ctx.LookupRoute(request)
// 	assert.True(t, ok)
// 	assert.NotNil(t, matched)

// 	// check it was cached
// 	_, ok = context.GetOk(request, ctxMatchedRoute)
// 	assert.True(t, ok)

// 	matched, ok = ctx.LookupRoute(request)
// 	assert.True(t, ok)
// 	assert.NotNil(t, matched)
// }

// func TestContextInvalidRoute(t *testing.T) {
// 	spec, api := petstore.NewAPI(t)
// 	ctx := NewContext(spec, api, nil)
// 	ctx.router = DefaultRouter(spec, ctx.api)

// 	request, _ := http.NewRequest("DELETE", "pets", nil)

// 	// check there's nothing there
// 	_, ok := context.GetOk(request, ctxMatchedRoute)
// 	assert.False(t, ok)

// 	matched, ok := ctx.LookupRoute(request)
// 	assert.False(t, ok)
// 	assert.Nil(t, matched)

// 	// check it was cached
// 	_, ok = context.GetOk(request, ctxMatchedRoute)
// 	assert.False(t, ok)

// 	matched, ok = ctx.LookupRoute(request)
// 	assert.False(t, ok)
// 	assert.Nil(t, matched)
// }

// func TestContextValidContentType(t *testing.T) {
// 	ct := "application/json"
// 	ctx := NewContext(nil, nil, nil)

// 	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
// 	request.Header.Set(httpkit.HeaderContentType, ct)

// 	// check there's nothing there
// 	_, ok := context.GetOk(request, ctxContentType)
// 	assert.False(t, ok)

// 	// trigger the parse
// 	mt, _, err := ctx.ContentType(request)
// 	assert.NoError(t, err)
// 	assert.Equal(t, ct, mt)

// 	// check it was cached
// 	_, ok = context.GetOk(request, ctxContentType)
// 	assert.True(t, ok)

// 	// check if the cast works and fetch from cache too
// 	mt, _, err = ctx.ContentType(request)
// 	assert.NoError(t, err)
// 	assert.Equal(t, ct, mt)
// }

// func TestContextInvalidContentType(t *testing.T) {
// 	ct := "application("
// 	ctx := NewContext(nil, nil, nil)

// 	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
// 	request.Header.Set(httpkit.HeaderContentType, ct)

// 	// check there's nothing there
// 	_, ok := context.GetOk(request, ctxContentType)
// 	assert.False(t, ok)

// 	// trigger the parse
// 	mt, _, err := ctx.ContentType(request)
// 	assert.Error(t, err)
// 	assert.Empty(t, mt)

// 	// check it was not cached
// 	_, ok = context.GetOk(request, ctxContentType)
// 	assert.False(t, ok)

// 	// check if the failure continues
// 	_, _, err = ctx.ContentType(request)
// 	assert.Error(t, err)
// }
