package bourbon

import (
	"encoding/json"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterAdd(t *testing.T) {
	r := createRouter()
	r.Add(&Route{Parent: new(bourbon), Method: "GET", Pattern: "/"})
	r.Add(&Route{Parent: new(bourbon), Method: "POST", Pattern: "/"})

	assert.Equal(t, 1, len(r.routes["GET"]))
	assert.Equal(t, 1, len(r.routes["POST"]))
}

func TestRouterFindRoute(t *testing.T) {
	r, req, rw := createRouterWithRoute()
	action := r.Find("OPTIONS", "/")
	action.Run(rw, req)
	assert.Equal(t, 200, rw.Code)
}

func TestRouterFindOptions(t *testing.T) {
	r, req, rw := createRouterWithRoute()
	action := r.Find("OPTIONS", "/")
	action.Run(rw, req)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, "0", rw.HeaderMap.Get("Content-Length"))
	assert.Equal(t, "GET", rw.HeaderMap.Get("Allow"))
}

func TestRouterFindPageNotFound(t *testing.T) {
	r, req, rw := createRouterWithRoute()
	action := r.Find("GET", "/404")
	action.Run(rw, req)

	msg, _ := json.Marshal(Message{Code: 404, Message: http.StatusText(404)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func TestRouterFindMethodNotAllowed(t *testing.T) {
	r, req, rw := createRouterWithRoute()
	action := r.Find("POST", "/")
	action.Run(rw, req)

	msg, _ := json.Marshal(Message{Code: 405, Message: http.StatusText(405)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 405, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func createRouter() *defaultRouter {
	return &defaultRouter{routes: make(map[string][]*Route)}
}

func createRouterWithRoute() (Router, *http.Request, *httptest.ResponseRecorder) {
	fn := func() int { return 200 }
	route := &Route{Parent: new(bourbon), Method: "GET", Pattern: "/", Handler: fn}
	r := createRouter()
	r.Add(route)
	req, _ := http.NewRequest("GET", "http://example.com/", strings.NewReader(""))
	rw := httptest.NewRecorder()
	return r, req, rw
}
