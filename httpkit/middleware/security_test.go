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
	"net/http/httptest"
	"testing"

	"github.com/go-swagger/go-swagger/internal/testing/petstore"
	"github.com/stretchr/testify/assert"
	netContext "golang.org/x/net/context"
)

func TestSecurityMiddleware(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	context.router = DefaultRouter(spec, context.api)
	mw := newSecureAPI(context, HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/pets", nil)

	r, _ := context.LookupRoute(request)
	ctx := NewContextWithMatchedRoute(netContext.TODO(), r)
	mw.ServeHTTP(ctx, recorder, request)
	assert.Equal(t, 401, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)
	request.SetBasicAuth("admin", "wrong")

    r, _ = context.LookupRoute(request)
	ctx = NewContextWithMatchedRoute(netContext.TODO(), r)
	mw.ServeHTTP(ctx, recorder, request)
	assert.Equal(t, 401, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)
	request.SetBasicAuth("admin", "admin")

	r, _ = context.LookupRoute(request)
	ctx = NewContextWithMatchedRoute(netContext.TODO(), r)
	mw.ServeHTTP(ctx, recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets/1", nil)

	r, _ = context.LookupRoute(request)
	ctx = NewContextWithMatchedRoute(netContext.TODO(), r)
	mw.ServeHTTP(ctx, recorder, request)
	assert.Equal(t, 200, recorder.Code)

}
