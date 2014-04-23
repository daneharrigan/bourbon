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
	rt := createDefaultRouter().(*defaultRouter)
	rt.Add(&route{parent: new(bourbon), method: "GET", pattern: "/"})
	rt.Add(&route{parent: new(bourbon), method: "POST", pattern: "/"})

	assert.Equal(t, 1, len(rt.routes["GET"]))
	assert.Equal(t, 1, len(rt.routes["POST"]))
}

func TestRouterFindRoute(t *testing.T) {
	rt, req, rw := createRouterWithRoute()
	r := rt.Find("OPTIONS", "/")
	c := createContext(r, rw, req)
	c.Run()
	assert.Equal(t, 200, rw.Code)
}

func TestRouterFindOptions(t *testing.T) {
	rt, req, rw := createRouterWithRoute()
	r := rt.Find("OPTIONS", "/")
	c := createContext(r, rw, req)
	c.Run()

	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, "0", rw.HeaderMap.Get("Content-Length"))
	assert.Equal(t, "GET", rw.HeaderMap.Get("Allow"))
}

func TestRouterFindPageNotFound(t *testing.T) {
	rt, req, rw := createRouterWithRoute()
	r := rt.Find("GET", "/404")
	c := createContext(r, rw, req)
	c.Run()

	msg, _ := json.Marshal(Message{Code: 404, Message: http.StatusText(404)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 404, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func TestRouterFindMethodNotAllowed(t *testing.T) {
	rt, req, rw := createRouterWithRoute()
	r := rt.Find("POST", "/")
	c := createContext(r, rw, req)
	c.Run()

	msg, _ := json.Marshal(Message{Code: 405, Message: http.StatusText(405)})
	enc, _ := ioutil.ReadAll(rw.Body)

	assert.Equal(t, 405, rw.Code)
	assert.Equal(t, msg, enc[:len(enc)-1])
}

func createRouterWithRoute() (Router, *http.Request, *httptest.ResponseRecorder) {
	fn := func() int { return 200 }
	rt := createDefaultRouter()
	rt.Add(&route{parent: new(bourbon), method: "GET", pattern: "/", handler: fn})
	r, _ := http.NewRequest("GET", "http://example.com/", strings.NewReader(""))
	rw := httptest.NewRecorder()
	return rt, r, rw
}
