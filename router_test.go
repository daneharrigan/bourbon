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
	rt := createRouter()
	rt.Add(&route{parent: new(bourbon), method: "GET", pattern: "/"})
	rt.Add(&route{parent: new(bourbon), method: "POST", pattern: "/"})

	assert.Equal(t, 1, len(rt.routes["GET"]))
	assert.Equal(t, 1, len(rt.routes["POST"]))
}

func TestRouterFindRoute(t *testing.T) {
	rt, r, rw := createRouterWithRoute()
	action := rt.Find("OPTIONS", "/")
	action.Run(rw, r)
	assert.Equal(t, 200, rw.Code)
}

func TestRouterFindOptions(t *testing.T) {
	rt, r, rw := createRouterWithRoute()
	action := rt.Find("OPTIONS", "/")
	action.Run(rw, r)

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, "0", rw.HeaderMap.Get("Content-Length"))
	assert.Equal(t, "GET", rw.HeaderMap.Get("Allow"))
}

func TestRouterFindPageNotFound(t *testing.T) {
	rt, r, rw := createRouterWithRoute()
	action := rt.Find("GET", "/404")
	action.Run(rw, r)

	msg, _ := json.Marshal(Message{Code: 404, Message: http.StatusText(404)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func TestRouterFindMethodNotAllowed(t *testing.T) {
	rt, r, rw := createRouterWithRoute()
	action := rt.Find("POST", "/")
	action.Run(rw, r)

	msg, _ := json.Marshal(Message{Code: 405, Message: http.StatusText(405)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 405, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func createRouter() *router {
	return &router{routes: make(map[string][]Route)}
}

func createRouterWithRoute() (Router, *http.Request, *httptest.ResponseRecorder) {
	fn := func() int { return 200 }
	rt := createRouter()
	rt.Add(&route{parent: new(bourbon), method: "GET", pattern: "/", handler: fn})
	r, _ := http.NewRequest("GET", "http://example.com/", strings.NewReader(""))
	rw := httptest.NewRecorder()
	return rt, r, rw
}
